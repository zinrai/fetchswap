package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	body, err := fetch(ctx, cfg.url)
	if err != nil {
		return err
	}
	defer body.Close()

	return swap(cfg.dest, body, cfg.mode)
}
