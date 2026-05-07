package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/watcher"
)

type watchConfig struct {
	left     string
	right    string
	format   string
	interval time.Duration
}

func parseWatchFlags(args []string) (watchConfig, error) {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text, json, csv")
	interval := fs.Duration("interval", 2*time.Second, "polling interval (e.g. 1s, 500ms)")

	if err := fs.Parse(args); err != nil {
		return watchConfig{}, err
	}
	if fs.NArg() < 2 {
		return watchConfig{}, fmt.Errorf("watch requires two .env file paths")
	}
	return watchConfig{
		left:     fs.Arg(0),
		right:    fs.Arg(1),
		format:   *format,
		interval: *interval,
	}, nil
}

func runWatch(args []string, out io.Writer) error {
	cfg, err := parseWatchFlags(args)
	if err != nil {
		return err
	}

	paths := []string{cfg.left, cfg.right}
	w, err := watcher.New(paths)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	fmt.Fprintf(out, "watching %s and %s (interval: %s)\n", cfg.left, cfg.right, cfg.interval)

	runDiff := func() error {
		left, err := parser.ParseFile(cfg.left)
		if err != nil {
			return err
		}
		right, err := parser.ParseFile(cfg.right)
		if err != nil {
			return err
		}
		result := diff.Compare(left, right)
		return formatter.Write(result, cfg.format, out)
	}

	if err := runDiff(); err != nil {
		return err
	}

	ticker := time.NewTicker(cfg.interval)
	defer ticker.Stop()
	for range ticker.C {
		events := w.Poll()
		for _, e := range events {
			if e.Err != nil {
				fmt.Fprintf(out, "error polling %s: %v\n", e.Path, e.Err)
				continue
			}
			if e.Changed {
				fmt.Fprintf(out, "--- change detected in %s ---\n", e.Path)
				if err := runDiff(); err != nil {
					fmt.Fprintf(out, "diff error: %v\n", err)
				}
				break
			}
		}
	}
	return nil
}
