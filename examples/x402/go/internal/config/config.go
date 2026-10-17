package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL        string `toml:"base_url"`
	Mode           string `toml:"mode"`
	Endpoint       string `toml:"endpoint"`
	Model          string `toml:"model"`
	Prompt         string `toml:"prompt"`
	N              int    `toml:"n"`
	ImagePath      string `toml:"image_path"`
	TimeoutSeconds int    `toml:"timeout_seconds"`
	EVMPrivateKey  string `toml:"evm_private_key"`
	EVMRPCURL      string `toml:"evm_rpc_url"`
}

func Load(path string) (Config, error) {
	_ = godotenv.Load(".env")
	cfg := Config{BaseURL: "https://api.ai-router.dev", Mode: "probe", Endpoint: "generations", Model: "gpt-image-2", Prompt: "A small red fox under a moonlit sky", N: 1, TimeoutSeconds: 60, EVMRPCURL: "https://mainnet.base.org"}
	if path == "" {
		path = os.Getenv("AI_ROUTER_X402_CONFIG")
	}
	if path == "" {
		path = "config.toml"
	}
	if data, err := os.ReadFile(filepath.Clean(path)); err == nil {
		if _, err := toml.Decode(string(data), &cfg); err != nil {
			return Config{}, fmt.Errorf("decode config.toml: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	cfg.BaseURL = firstEnv("AI_ROUTER_X402_BASE_URL", cfg.BaseURL)
	cfg.Mode = firstEnv("AI_ROUTER_X402_MODE", cfg.Mode)
	cfg.Endpoint = firstEnv("AI_ROUTER_X402_ENDPOINT", cfg.Endpoint)
	cfg.Model = firstEnv("AI_ROUTER_X402_MODEL", cfg.Model)
	cfg.Prompt = firstEnv("AI_ROUTER_X402_PROMPT", cfg.Prompt)
	cfg.ImagePath = firstEnv("AI_ROUTER_X402_IMAGE_PATH", cfg.ImagePath)
	cfg.EVMPrivateKey = firstEnv("AI_ROUTER_X402_EVM_PRIVATE_KEY", cfg.EVMPrivateKey)
	cfg.EVMRPCURL = firstEnv("AI_ROUTER_X402_EVM_RPC_URL", cfg.EVMRPCURL)
	if raw := os.Getenv("AI_ROUTER_X402_N"); raw != "" {
		cfg.N, _ = strconv.Atoi(raw)
	}
	if raw := os.Getenv("AI_ROUTER_X402_TIMEOUT_SECONDS"); raw != "" {
		cfg.TimeoutSeconds, _ = strconv.Atoi(raw)
	}
	cfg.BaseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	cfg.Mode = strings.ToLower(strings.TrimSpace(cfg.Mode))
	cfg.Endpoint = strings.ToLower(strings.TrimSpace(cfg.Endpoint))
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
func (c Config) Timeout() time.Duration { return time.Duration(c.TimeoutSeconds) * time.Second }
func (c Config) Validate() error {
	if c.Mode != "probe" && c.Mode != "pay" {
		return errors.New("mode must be probe or pay")
	}
	if c.Endpoint != "generations" && c.Endpoint != "edits" {
		return errors.New("endpoint must be generations or edits")
	}
	if c.Model == "" || c.Prompt == "" || c.N < 1 || c.N > 4 || c.TimeoutSeconds <= 0 {
		return errors.New("model, prompt, n (1-4), and positive timeout are required")
	}
	if c.Mode == "pay" && !strings.HasPrefix(strings.TrimSpace(c.EVMPrivateKey), "0x") {
		return errors.New("pay mode requires EVM_PRIVATE_KEY in .env or config.toml")
	}
	if c.Endpoint == "edits" && c.Mode == "pay" && strings.TrimSpace(c.ImagePath) == "" {
		return errors.New("edits pay mode requires AI_ROUTER_X402_IMAGE_PATH")
	}
	return nil
}
func firstEnv(name, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	return fallback
}
