// Package config contains configuration shared by the CLI and API client.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is deliberately small: transport concerns stay in the client while
// the CLI decides where to read the values from (flags override these values).
type Config struct {
	APIKey         string
	BaseURL        string
	Model          string
	Size           string
	ResponseFormat string
	Timeout        time.Duration
}

// FromEnv returns production-safe defaults for an OpenAI-compatible Images API.
// OPENAI_* names are accepted first; AI_ROUTER_* aliases make the example easy
// to use with AI-ROUTER without changing application code.
func FromEnv() Config {
	timeout := 120 * time.Second
	if raw := firstNonEmpty(os.Getenv("IMAGE_API_TIMEOUT"), os.Getenv("AI_ROUTER_IMAGE_TIMEOUT")); raw != "" {
		if seconds, err := strconv.Atoi(raw); err == nil && seconds > 0 {
			timeout = time.Duration(seconds) * time.Second
		}
	}

	return Config{
		APIKey:         firstNonEmpty(os.Getenv("OPENAI_API_KEY"), os.Getenv("AI_ROUTER_API_KEY")),
		BaseURL:        strings.TrimRight(firstNonEmpty(os.Getenv("OPENAI_BASE_URL"), os.Getenv("AI_ROUTER_BASE_URL"), "https://api.ai-router.dev/v1"), "/"),
		Model:          firstNonEmpty(os.Getenv("IMAGE_MODEL"), os.Getenv("AI_ROUTER_MODEL"), "gpt-image-2"),
		Size:           firstNonEmpty(os.Getenv("IMAGE_SIZE"), os.Getenv("AI_ROUTER_SIZE"), "1024x1024"),
		ResponseFormat: firstNonEmpty(os.Getenv("IMAGE_RESPONSE_FORMAT"), "b64_json"),
		Timeout:        timeout,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
