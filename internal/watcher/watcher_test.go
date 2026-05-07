package watcher_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestNew_ValidFile(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	_, err := watcher.New([]string{p})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := watcher.New([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestPoll_NoChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	w, err := watcher.New([]string{p})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	events := w.Poll()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Changed {
		t.Error("expected no change")
	}
}

func TestPoll_DetectsChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=old\n")
	w, err := watcher.New([]string{p})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := os.WriteFile(p, []byte("KEY=new\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	events := w.Poll()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if !events[0].Changed {
		t.Error("expected change to be detected")
	}
}

func TestPoll_SubsequentPollNoChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=old\n")
	w, _ := watcher.New([]string{p})
	os.WriteFile(p, []byte("KEY=new\n"), 0644)
	w.Poll() // consume the change
	events := w.Poll()
	if events[0].Changed {
		t.Error("second poll should not report change")
	}
}

func TestPoll_ErrorOnDeletedFile(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	w, _ := watcher.New([]string{p})
	os.Remove(p)
	events := w.Poll()
	if events[0].Err == nil {
		t.Error("expected error for deleted file")
	}
}
