package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

// exportFlags holds CLI flags for the export sub-command.
type exportFlags struct {
	format string
	prefix string
	file   string
}

func parseExportFlags(args []string) (exportFlags, error) {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	var f exportFlags
	fs.StringVar(&f.format, "format", "dotenv", "output format: dotenv, shell, json")
	fs.StringVar(&f.prefix, "prefix", "", "only export keys with this prefix")
	if err := fs.Parse(args); err != nil {
		return f, err
	}
	if fs.NArg() < 1 {
		return f, fmt.Errorf("usage: envdiff export [flags] <file>")
	}
	f.file = fs.Arg(0)
	return f, nil
}

// runExport parses a .env file and writes it in the requested format.
func runExport(args []string) error {
	flags, err := parseExportFlags(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(flags.file)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", flags.file, err)
	}

	opts := exporter.Options{
		Format: exporter.Format(flags.format),
		Prefix: flags.prefix,
	}

	return exporter.Export(os.Stdout, env, opts)
}
