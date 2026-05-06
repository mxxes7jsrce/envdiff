package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// Config holds the parsed CLI configuration.
type Config struct {
	LeftFile  string
	RightFile string
	Quiet     bool
}

// parseFlags parses command-line arguments and returns a Config.
// It writes usage information to stderr on error.
func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	quiet := fs.Bool("quiet", false, "suppress output; exit code only")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <file1> <file2>\n\nFlags:\n")
		fs.SetOutput(os.Stderr)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fs.Usage()
		return nil, err
	}

	if fs.NArg() != 2 {
		fs.Usage()
		return nil, errors.New("exactly two .env files are required")
	}

	return &Config{
		LeftFile:  fs.Arg(0),
		RightFile: fs.Arg(1),
		Quiet:     *quiet,
	}, nil
}
