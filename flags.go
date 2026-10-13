package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type config struct {
	url     string
	dest    string
	mode    os.FileMode
	timeout time.Duration

	showVersion bool
}

func parseFlags(args []string) (config, error) {
	var cfg config
	var mode string
	fs := flag.NewFlagSet("fetchswap", flag.ExitOnError)
	fs.StringVar(&cfg.url, "url", "", "URL to fetch")
	fs.StringVar(&cfg.dest, "dest", "", "file to replace")
	fs.StringVar(&mode, "mode", "", "permission bits of dest in octal, such as 0644")
	fs.DurationVar(&cfg.timeout, "timeout", 0, "limit for the whole request, such as 60s")
	fs.BoolVar(&cfg.showVersion, "version", false, "print the version and exit")
	fs.Parse(args)

	if cfg.showVersion {
		return cfg, nil
	}

	if err := checkRequired(fs, "url", "dest", "mode", "timeout"); err != nil {
		return config{}, err
	}
	m, err := parseMode(mode)
	if err != nil {
		return config{}, err
	}
	cfg.mode = m
	return cfg, nil
}

// No defaults: the mode a file needs and how long a job may wait are the
// caller's decisions.
func checkRequired(fs *flag.FlagSet, names ...string) error {
	// No positional arguments: an unquoted space in -dest would otherwise be
	// dropped silently and the file written elsewhere.
	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	for _, name := range names {
		if !set[name] {
			return fmt.Errorf("missing -%s", name)
		}
	}
	return nil
}

func parseMode(s string) (os.FileMode, error) {
	v, err := strconv.ParseUint(s, 8, 32)
	// Reject above 0777, not mask: Chmod drops higher bits silently, so 1777
	// would become 0777.
	if err != nil || v > 0o777 {
		return 0, fmt.Errorf("invalid mode %q: must be octal between 0000 and 0777", s)
	}
	return os.FileMode(v), nil
}
