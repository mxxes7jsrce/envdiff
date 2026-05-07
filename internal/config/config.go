// Package config handles loading and validating envdiff runtime configuration
// from environment variables and optional config files.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the runtime configuration for envdiff.
type Config struct {
	// DefaultFormat specifies the output format when none is provided via flag.
	// Accepted values: "text", "json", "csv". Defaults to "text".
	DefaultFormat string

	// DefaultPrefix is a key prefix filter applied when none is provided via flag.
	DefaultPrefix string

	// IgnoreKeys is a comma-separated list of keys to always exclude from output.
	IgnoreKeys []string

	// Quiet suppresses non-error output when set to true.
	Quiet bool
}

const (
	envDefaultFormat = "ENVDIFF_FORMAT"
	envDefaultPrefix = "ENVDIFF_PREFIX"
	envIgnoreKeys    = "ENVDIFF_IGNORE_KEYS"
	envQuiet         = "ENVDIFF_QUIET"
)

var validFormats = map[string]bool{
	"text": true,
	"json": true,
	"csv":  true,
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if any value fails validation.
func Load() (*Config, error) {
	cfg := &Config{
		DefaultFormat: "text",
	}

	if v := os.Getenv(envDefaultFormat); v != "" {
		norm := strings.ToLower(strings.TrimSpace(v))
		if !validFormats[norm] {
			return nil, fmt.Errorf("config: invalid %s value %q: must be one of text, json, csv", envDefaultFormat, v)
		}
		cfg.DefaultFormat = norm
	}

	if v := os.Getenv(envDefaultPrefix); v != "" {
		cfg.DefaultPrefix = strings.TrimSpace(v)
	}

	if v := os.Getenv(envIgnoreKeys); v != "" {
		parts := strings.Split(v, ",")
		for _, p := range parts {
			key := strings.TrimSpace(p)
			if key != "" {
				cfg.IgnoreKeys = append(cfg.IgnoreKeys, key)
			}
		}
	}

	if v := os.Getenv(envQuiet); v != "" {
		norm := strings.ToLower(strings.TrimSpace(v))
		cfg.Quiet = norm == "1" || norm == "true" || norm == "yes"
	}

	return cfg, nil
}
