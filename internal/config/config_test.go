package config_test

import (
	"testing"

	"github.com/user/envdiff/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("ENVDIFF_FORMAT", "")
	t.Setenv("ENVDIFF_PREFIX", "")
	t.Setenv("ENVDIFF_IGNORE_KEYS", "")
	t.Setenv("ENVDIFF_QUIET", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultFormat != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.DefaultFormat)
	}
	if cfg.Quiet {
		t.Error("expected Quiet to be false by default")
	}
	if len(cfg.IgnoreKeys) != 0 {
		t.Errorf("expected no ignore keys, got %v", cfg.IgnoreKeys)
	}
}

func TestLoad_CustomFormat(t *testing.T) {
	for _, fmt := range []string{"json", "csv", "text", "JSON"} {
		t.Run(fmt, func(t *testing.T) {
			t.Setenv("ENVDIFF_FORMAT", fmt)
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("unexpected error for format %q: %v", fmt, err)
			}
			if cfg.DefaultFormat == "" {
				t.Errorf("expected non-empty format for input %q", fmt)
			}
		})
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	t.Setenv("ENVDIFF_FORMAT", "yaml")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestLoad_IgnoreKeys(t *testing.T) {
	t.Setenv("ENVDIFF_IGNORE_KEYS", "SECRET, TOKEN , DEBUG")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.IgnoreKeys) != 3 {
		t.Fatalf("expected 3 ignore keys, got %d: %v", len(cfg.IgnoreKeys), cfg.IgnoreKeys)
	}
	if cfg.IgnoreKeys[1] != "TOKEN" {
		t.Errorf("expected trimmed key 'TOKEN', got %q", cfg.IgnoreKeys[1])
	}
}

func TestLoad_QuietVariants(t *testing.T) {
	for _, val := range []string{"1", "true", "yes", "TRUE", "YES"} {
		t.Run(val, func(t *testing.T) {
			t.Setenv("ENVDIFF_QUIET", val)
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !cfg.Quiet {
				t.Errorf("expected Quiet=true for value %q", val)
			}
		})
	}
}

func TestLoad_Prefix(t *testing.T) {
	t.Setenv("ENVDIFF_PREFIX", "APP_")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultPrefix != "APP_" {
		t.Errorf("expected prefix 'APP_', got %q", cfg.DefaultPrefix)
	}
}
