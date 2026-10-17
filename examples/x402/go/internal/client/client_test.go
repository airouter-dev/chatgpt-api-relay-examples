package client

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDecodeChallengeValidatesPaymentTerms(t *testing.T) {
	payload := `{"resource":{"url":"https://api.ai-router.dev/x402/v1/images/generations"},"accepts":[{"scheme":"exact","network":"eip155:8453","asset":"0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913","amount":"4500","payTo":"0xEe027fE58800Bcef81469EDF0C2E7f2ab54965De"}]}`
	response := &http.Response{StatusCode: http.StatusPaymentRequired, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(""))}
	response.Header.Set("PAYMENT-REQUIRED", base64.RawURLEncoding.EncodeToString([]byte(payload)))
	challenge, err := decodeChallenge(response)
	if err != nil || challenge.Amount != "4500" || challenge.Network != network {
		t.Fatalf("unexpected valid challenge: %#v, %v", challenge, err)
	}
	bad := strings.Replace(payload, payTo, "0x0000000000000000000000000000000000000001", 1)
	response.Header.Set("PAYMENT-REQUIRED", base64.RawURLEncoding.EncodeToString([]byte(bad)))
	if _, err := decodeChallenge(response); err == nil {
		t.Fatal("expected recipient validation failure")
	}
}
