// Command imagegen demonstrates a complete text-to-image then image-to-image
// workflow against any OpenAI-compatible Images API.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ai-router/image-generation-multilang/go/internal/client"
	"github.com/ai-router/image-generation-multilang/go/internal/config"
)

type flags struct {
	prompt       string
	editPrompt   string
	output       string
	editedOutput string
	apiKey       string
	baseURL      string
	model        string
	size         string
	format       string
	timeout      int
}

func main() {
	cfg := config.FromEnv()
	options := parseFlags(cfg)
	cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.Size, cfg.ResponseFormat = options.apiKey, options.baseURL, options.model, options.size, options.format
	cfg.Timeout = time.Duration(options.timeout) * time.Second

	if strings.TrimSpace(options.prompt) == "" || strings.TrimSpace(options.editPrompt) == "" {
		fatal(errors.New("both --prompt and --edit-prompt are required"))
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		fatal(errors.New("API key is missing; pass --api-key or set OPENAI_API_KEY"))
	}

	api := client.New(cfg)
	ctx := context.Background()
	generated, err := api.Generate(ctx, options.prompt)
	if err != nil {
		fatal(fmt.Errorf("text-to-image: %w", err))
	}
	if err := writeImage(options.output, generated); err != nil {
		fatal(fmt.Errorf("save generated image: %w", err))
	}

	edited, err := api.Edit(ctx, options.output, options.editPrompt)
	if err != nil {
		fatal(fmt.Errorf("image-to-image: %w", err))
	}
	if err := writeImage(options.editedOutput, edited); err != nil {
		fatal(fmt.Errorf("save edited image: %w", err))
	}
	fmt.Printf("Generated %s\nEdited    %s\n", options.output, options.editedOutput)
}

func parseFlags(cfg config.Config) flags {
	result := flags{}
	flag.StringVar(&result.prompt, "prompt", "", "text prompt for the first image (required)")
	flag.StringVar(&result.editPrompt, "edit-prompt", "", "instruction applied by /images/edits (required)")
	flag.StringVar(&result.output, "output", "generated.png", "path for the generated image")
	flag.StringVar(&result.editedOutput, "edited-output", "edited.png", "path for the edited image")
	flag.StringVar(&result.apiKey, "api-key", cfg.APIKey, "API key (prefer OPENAI_API_KEY/AI_ROUTER_API_KEY)")
	flag.StringVar(&result.baseURL, "base-url", cfg.BaseURL, "API base URL, e.g. https://api.ai-router.dev/v1")
	flag.StringVar(&result.model, "model", cfg.Model, "image model: gpt-image-2, gpt-image-2.5, ...")
	flag.StringVar(&result.size, "size", cfg.Size, "image size accepted by the selected model")
	flag.StringVar(&result.format, "response-format", cfg.ResponseFormat, "b64_json or url")
	flag.IntVar(&result.timeout, "timeout", int(cfg.Timeout/time.Second), "HTTP timeout in seconds")
	flag.Parse()
	return result
}

func writeImage(path string, data []byte) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("output path must not be empty")
	}
	if len(data) == 0 {
		return errors.New("image is empty")
	}
	if parent := filepath.Dir(path); parent != "." {
		if err := os.MkdirAll(parent, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0o644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
