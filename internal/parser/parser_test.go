package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", env["APP_ENV"], "production")
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q, want %q", env["DB_HOST"], "localhost")
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	content := "# this is a comment\n\nFOO=bar\n"
	path := writeTempEnv(t, content)

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestParseFile_StripQuotes(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret"\nTOKEN='abc123'\n`)

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret" {
		t.Errorf("SECRET: got %q, want %q", env["SECRET"], "my secret")
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseLine_MissingEquals(t *testing.T) {
	_, _, err := parseLine("NODIVIDER")
	if err == nil {
		t.Error("expected error for line without '=', got nil")
	}
}

func TestParseLine_EmptyKey(t *testing.T) {
	_, _, err := parseLine("=value")
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}
