package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents a supported export format.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// Options controls export behaviour.
type Options struct {
	Format  Format
	Prefix  string
	Redact  bool
}

// Export writes the given key-value map to w in the requested format.
func Export(w io.Writer, env map[string]string, opts Options) error {
	keys := sortedKeys(env)

	switch opts.Format {
	case FormatShell:
		return exportShell(w, env, keys, opts.Prefix)
	case FormatJSON:
		return exportJSON(w, env, keys)
	case FormatDotEnv, "":
		return exportDotEnv(w, env, keys, opts.Prefix)
	default:
		return fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func exportDotEnv(w io.Writer, env map[string]string, keys []string, prefix string) error {
	for _, k := range keys {
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		if _, err := fmt.Fprintf(w, "%s=%q\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

func exportShell(w io.Writer, env map[string]string, keys []string, prefix string) error {
	for _, k := range keys {
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		if _, err := fmt.Fprintf(w, "export %s=%q\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(w io.Writer, env map[string]string, keys []string) error {
	ordered := make([]map[string]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, map[string]string{"key": k, "value": env[k]})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

// ParseFormat converts a string to a known Format, returning an error if the
// value is not recognised. An empty string resolves to FormatDotEnv.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatDotEnv, "":
		return FormatDotEnv, nil
	case FormatShell:
		return FormatShell, nil
	case FormatJSON:
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", s)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
