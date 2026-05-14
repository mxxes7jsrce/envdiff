package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/baseline"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	env := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"LOG_LEVEL": "info",
	}

	if err := baseline.Save(path, ".env.production", env); err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if snap.File != ".env.production" {
		t.Errorf("File = %q, want %q", snap.File, ".env.production")
	}

	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if snap.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt should not be in the future")
	}

	for k, want := range env {
		got, ok := snap.Env[k]
		if !ok {
			t.Errorf("missing key %q in loaded snapshot", k)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("Load() expected error for missing file, got nil")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := baseline.Save("/nonexistent/dir/baseline.json", "src", map[string]string{"K": "V"})
	if err == nil {
		t.Fatal("Save() expected error for invalid path, got nil")
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")

	if err := os.WriteFile(path, []byte("not valid json{"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("Load() expected error for corrupt JSON, got nil")
	}
}
