// Command envdiff compares .env files across environments and surfaces
// missing or mismatched keys.
package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(cfg.LeftFile)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", cfg.LeftFile, err)
	}

	right, err := parser.ParseFile(cfg.RightFile)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", cfg.RightFile, err)
	}

	result := diff.Compare(left, right)

	w := os.Stdout
	reporter.Report(w, result, cfg.LeftFile, cfg.RightFile)

	if !result.Clean() {
		os.Exit(2)
	}
	return nil
}
