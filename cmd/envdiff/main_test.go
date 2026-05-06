package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp env: %v", err)
	}
	return path
}

func TestRun_Clean(t *testing.T) {
	left := writeTempEnv(t, "KEY=value\nDEBUG=true\n")
	right := writeTempEnv(t, "KEY=value\nDEBUG=true\n")

	if err := run([]string{left, right}); err != nil {
		t.Errorf("expected no error for identical files, got: %v", err)
	}
}

func TestRun_MissingFile(t *testing.T) {
	left := writeTempEnv(t, "KEY=value\n")

	if err := run([]string{left, "/nonexistent/.env"}); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRun_NoArgs(t *testing.T) {
	if err := run([]string{}); err == nil {
		t.Error("expected error when no arguments provided")
	}
}
