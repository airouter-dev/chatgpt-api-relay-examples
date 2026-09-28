// Package client implements the narrow Images API surface used by this demo.
// It does not expose HTTP details to the CLI, which keeps the generation/edit
// workflow easy to reuse from another Go service.
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
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ai-router/image-generation-multilang/go/internal/config"
)

const maxResponseBytes = 100 << 20 // protect a CLI process from an accidental huge response

// Client is safe to reuse for multiple requests. HTTPClient may be replaced in
// tests or by callers that need a custom transport.
type Client struct {
	Config     config.Config
	HTTPClient *http.Client
}

// New creates a client with a timeout even when a caller forgot to configure
// one. The API key is checked lazily so tests can exercise a local server.
func New(cfg config.Config) *Client {
	transportTimeout := cfg.Timeout
	if transportTimeout <= 0 {
		transportTimeout = 120 * time.Second
	}
	return &Client{Config: cfg, HTTPClient: &http.Client{Timeout: transportTimeout}}
}

type imageResponse struct {
	Data []struct {
		B64JSON string `json:"b64_json"`
		URL     string `json:"url"`
	} `json:"data"`
}

type apiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Generate requests a new image and resolves either b64_json or url into bytes.
func (c *Client) Generate(ctx context.Context, prompt string) ([]byte, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, errors.New("prompt must not be empty")
	}
	body := map[string]string{
		"model":           c.Config.Model,
		"prompt":          prompt,
		"size":            c.Config.Size,
		"response_format": c.Config.ResponseFormat,
	}
	return c.postJSON(ctx, "/images/generations", body)
}

// Edit sends an image generated (or supplied by the caller) to /images/edits.
// The image field follows the OpenAI-compatible multipart contract.
func (c *Client) Edit(ctx context.Context, imagePath, prompt string) ([]byte, error) {
	if strings.TrimSpace(imagePath) == "" {
		return nil, errors.New("image path must not be empty")
	}
	if strings.TrimSpace(prompt) == "" {
		return nil, errors.New("edit prompt must not be empty")
	}

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("open edit image %q: %w", imagePath, err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fields := map[string]string{
		"model":           c.Config.Model,
		"prompt":          prompt,
		"size":            c.Config.Size,
		"response_format": c.Config.ResponseFormat,
	}
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("write multipart field %s: %w", key, err)
		}
	}
	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return nil, fmt.Errorf("create image multipart field: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("read edit image: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint("/images/edits"), &body)
	if err != nil {
		return nil, fmt.Errorf("create edit request: %w", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return c.resolveResponse(c.do(request))
}

func (c *Client) postJSON(ctx context.Context, path string, payload map[string]string) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(path), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	return c.resolveResponse(c.do(request))
}

func (c *Client) resolveResponse(response *http.Response, requestErr error) ([]byte, error) {
	if requestErr != nil {
		return nil, requestErr
	}
	defer response.Body.Close()
	limited := io.LimitReader(response.Body, maxResponseBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read API response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		var apiErr apiErrorResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("image API returned %s: %s", response.Status, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("image API returned %s: %s", response.Status, trimBody(body))
	}

	var decoded imageResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("decode image response: %w", err)
	}
	if len(decoded.Data) == 0 {
		return nil, errors.New("image API returned no data items")
	}
	item := decoded.Data[0]
	if item.B64JSON != "" {
		return decodeBase64Image(item.B64JSON)
	}
	if item.URL != "" {
		return c.download(item.URL)
	}
	return nil, errors.New("image API data item contains neither b64_json nor url")
}

func (c *Client) do(request *http.Request) (*http.Response, error) {
	if strings.TrimSpace(c.Config.APIKey) == "" {
		return nil, errors.New("API key is missing; set OPENAI_API_KEY or AI_ROUTER_API_KEY")
	}
	request.Header.Set("Authorization", "Bearer "+c.Config.APIKey)
	request.Header.Set("Accept", "application/json")
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	return response, nil
}

func (c *Client) download(rawURL string) ([]byte, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid image URL in API response: %q", rawURL)
	}
	request, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create image download request: %w", err)
	}
	request.Header.Set("Accept", "image/*")
	// Image URLs are returned by the provider and may use object storage. Do not
	// forward the API bearer token to a different host.
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("image download request failed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("image download returned %s", response.Status)
	}
	bytes, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("read image download: %w", err)
	}
	if len(bytes) == 0 {
		return nil, errors.New("image download was empty")
	}
	return bytes, nil
}

func (c *Client) endpoint(path string) string {
	return strings.TrimRight(c.Config.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func decodeBase64Image(raw string) ([]byte, error) {
	value := strings.TrimSpace(raw)
	if comma := strings.Index(value, ","); strings.HasPrefix(value, "data:") && comma >= 0 {
		value = value[comma+1:]
	}
	value = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, value)
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("decode image b64_json: %w", err)
	}
	if len(decoded) == 0 {
		return nil, errors.New("image b64_json decoded to empty bytes")
	}
	return decoded, nil
}

func trimBody(body []byte) string {
	const max = 512
	text := strings.TrimSpace(string(body))
	if len(text) > max {
		return text[:max] + "..."
	}
	return text
}
