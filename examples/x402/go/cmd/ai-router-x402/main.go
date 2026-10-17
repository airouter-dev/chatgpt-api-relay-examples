package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/airouter-dev/chatgpt-api-relay-examples/examples/x402/go/internal/client"
	"github.com/airouter-dev/chatgpt-api-relay-examples/examples/x402/go/internal/config"
	"os"
)

func main() {
	p := flag.String("config", "", "optional TOML file")
	m := flag.String("mode", "", "probe or pay")
	e := flag.String("endpoint", "", "generations or edits")
	flag.Parse()
	cfg, err := config.Load(*p)
	if err != nil {
		fatal(err)
	}
	if *m != "" {
		cfg.Mode = *m
	}
	if *e != "" {
		cfg.Endpoint = *e
	}
	if err = cfg.Validate(); err != nil {
		fatal(err)
	}
	_, err = client.Run(context.Background(), cfg)
	if err != nil {
		fatal(err)
	}
}
func fatal(err error) { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
