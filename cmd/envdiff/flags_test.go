package main

import (
	"testing"
)

func TestParseFlags_ValidArgs(t *testing.T) {
	cfg, err := parseFlags([]string{".env", ".env.production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LeftFile != ".env" {
		t.Errorf("LeftFile = %q, want %q", cfg.LeftFile, ".env")
	}
	if cfg.RightFile != ".env.production" {
		t.Errorf("RightFile = %q, want %q", cfg.RightFile, ".env.production")
	}
	if cfg.Quiet {
		t.Error("Quiet should default to false")
	}
}

func TestParseFlags_QuietFlag(t *testing.T) {
	cfg, err := parseFlags([]string{"-quiet", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Quiet {
		t.Error("expected Quiet to be true")
	}
}

func TestParseFlags_MissingArgs(t *testing.T) {
	_, err := parseFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing arguments")
	}
}

func TestParseFlags_TooManyArgs(t *testing.T) {
	_, err := parseFlags([]string{"a.env", "b.env", "c.env"})
	if err == nil {
		t.Fatal("expected error for too many arguments")
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "a.env", "b.env"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
