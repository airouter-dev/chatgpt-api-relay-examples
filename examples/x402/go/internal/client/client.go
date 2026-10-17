package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/airouter-dev/chatgpt-api-relay-examples/examples/x402/go/internal/config"
	x402 "github.com/x402-foundation/x402/go/v2"
	x402http "github.com/x402-foundation/x402/go/v2/http"
	evmexact "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/client"
	evmsigners "github.com/x402-foundation/x402/go/v2/signers/evm"
)

const network = "eip155:8453"
const asset = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
const payTo = "0xEe027fE58800Bcef81469EDF0C2E7f2ab54965De"

type Challenge struct {
	Scheme, Network, Asset, PayTo, Amount, URL string
	Decimals                                   int
}

func Run(ctx context.Context, cfg config.Config) (Challenge, error) {
	req, err := request(cfg)
	if err != nil {
		return Challenge{}, err
	}
	if cfg.Mode == "pay" {
		signer, err := evmsigners.NewClientSignerFromPrivateKey(cfg.EVMPrivateKey)
		if err != nil {
			return Challenge{}, fmt.Errorf("create EVM signer: %w", err)
		}
		pc := x402.Newx402Client(x402.WithPaymentSelector(func(requirements []x402.PaymentRequirementsView) x402.PaymentRequirementsView {
			for _, r := range requirements {
				if r.GetScheme() == "exact" && r.GetNetwork() == network && strings.EqualFold(r.GetAsset(), asset) && strings.EqualFold(r.GetPayTo(), payTo) {
					return r
				}
			}
			panic("402 challenge did not contain the expected exact Base USDC recipient")
		})).Register(network, evmexact.NewExactEvmScheme(signer, nil))
		hc := x402http.WrapHTTPClientWithPayment(&http.Client{Timeout: cfg.Timeout()}, x402http.Newx402HTTPClient(pc))
		resp, err := hc.Do(req.WithContext(ctx))
		if err != nil {
			return Challenge{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusPaymentRequired {
			return decodeChallenge(resp)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return Challenge{}, fmt.Errorf("paid request returned %s: %s", resp.Status, strings.TrimSpace(string(b)))
		}
		b, _ := io.ReadAll(resp.Body)
		fmt.Printf("paid response: %s\n", strings.TrimSpace(string(b)))
		return Challenge{}, nil
	}
	resp, err := (&http.Client{Timeout: cfg.Timeout()}).Do(req.WithContext(ctx))
	if err != nil {
		return Challenge{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPaymentRequired {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return Challenge{}, fmt.Errorf("probe returned %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	return decodeChallenge(resp)
}
func request(cfg config.Config) (*http.Request, error) {
	endpoint := cfg.BaseURL + "/x402/v1/images/" + cfg.Endpoint
	if cfg.Endpoint == "generations" {
		b, _ := json.Marshal(map[string]any{"model": cfg.Model, "prompt": cfg.Prompt, "n": cfg.N})
		r, e := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
		if e == nil {
			r.Header.Set("Content-Type", "application/json")
		}
		return r, e
	}
	f, e := os.Open(cfg.ImagePath)
	if e != nil {
		return nil, fmt.Errorf("open edit image: %w", e)
	}
	defer f.Close()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for n, v := range map[string]string{"model": cfg.Model, "prompt": cfg.Prompt, "n": fmt.Sprint(cfg.N)} {
		_ = w.WriteField(n, v)
	}
	p, e := w.CreateFormFile("image", filepath.Base(cfg.ImagePath))
	if e != nil {
		return nil, e
	}
	if _, e = io.Copy(p, f); e != nil {
		return nil, e
	}
	_ = w.Close()
	r, e := http.NewRequest(http.MethodPost, endpoint, &body)
	if e == nil {
		r.Header.Set("Content-Type", w.FormDataContentType())
	}
	return r, e
}
func decodeChallenge(resp *http.Response) (Challenge, error) {
	raw := resp.Header.Get("PAYMENT-REQUIRED")
	if raw == "" {
		return Challenge{}, errors.New("402 response has no PAYMENT-REQUIRED header")
	}
	decoded, e := base64.RawURLEncoding.DecodeString(raw)
	if e != nil {
		decoded, e = base64.URLEncoding.DecodeString(raw)
	}
	if e != nil {
		return Challenge{}, fmt.Errorf("decode PAYMENT-REQUIRED: %w", e)
	}
	var env struct {
		Resource struct {
			URL string `json:"url"`
		} `json:"resource"`
		Accepts []struct {
			Scheme  string `json:"scheme"`
			Network string `json:"network"`
			Asset   string `json:"asset"`
			Amount  string `json:"amount"`
			PayTo   string `json:"payTo"`
		} `json:"accepts"`
	}
	if e = json.Unmarshal(decoded, &env); e != nil || len(env.Accepts) == 0 {
		return Challenge{}, errors.New("invalid PAYMENT-REQUIRED payload")
	}
	a := env.Accepts[0]
	c := Challenge{a.Scheme, a.Network, a.Asset, a.PayTo, a.Amount, env.Resource.URL, 6}
	if c.Scheme != "exact" || c.Network != network || !strings.EqualFold(c.Asset, asset) || !strings.EqualFold(c.PayTo, payTo) || c.Amount == "" {
		return c, fmt.Errorf("unexpected payment terms: %+v", c)
	}
	fmt.Printf("402 challenge: scheme=%s network=%s amount=%s asset=%s payTo=%s\n", c.Scheme, c.Network, c.Amount, c.Asset, c.PayTo)
	return c, nil
}
