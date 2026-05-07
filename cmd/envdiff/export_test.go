package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeExportTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestRunExport_DotEnv(t *testing.T) {
	p := writeExportTempEnv(t, "APP=hello\nDB_HOST=localhost\n")
	if err := runExport([]string{p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunExport_ShellFormat(t *testing.T) {
	p := writeExportTempEnv(t, "APP=hello\n")
	if err := runExport([]string{"-format", "shell", p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunExport_JSONFormat(t *testing.T) {
	p := writeExportTempEnv(t, "APP=hello\n")
	if err := runExport([]string{"-format", "json", p}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunExport_MissingFile(t *testing.T) {
	err := runExport([]string{"/no/such/file.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunExport_NoArgs(t *testing.T) {
	err := runExport([]string{})
	if err == nil {
		t.Fatal("expected error when no file provided")
	}
}

func TestParseExportFlags_PrefixFlag(t *testing.T) {
	flags, err := parseExportFlags([]string{"-prefix", "DB_", "file.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.prefix != "DB_" {
		t.Errorf("expected prefix DB_, got %q", flags.prefix)
	}
	if flags.file != "file.env" {
		t.Errorf("expected file file.env, got %q", flags.file)
	}
}
