package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshotter"
)

func TestTake_CopiesEnv(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshotter.Take("test.env", env)
	env["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected snapshot to be isolated from source map")
	}
	if s.Source != "test.env" {
		t.Errorf("expected source %q, got %q", "test.env", s.Source)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"KEY": "value", "PORT": "8080"}
	s := snapshotter.Take("prod.env", env)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshotter.SaveTo(tmp, s); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}
	loaded, err := snapshotter.LoadFrom(tmp)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}
	if loaded.Source != s.Source {
		t.Errorf("source mismatch: got %q want %q", loaded.Source, s.Source)
	}
	for k, v := range env {
		if loaded.Env[k] != v {
			t.Errorf("key %q: got %q want %q", k, loaded.Env[k], v)
		}
	}
}

func TestLoadFrom_MissingFile(t *testing.T) {
	_, err := snapshotter.LoadFrom("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveTo_InvalidPath(t *testing.T) {
	s := snapshotter.Take("x", map[string]string{})
	err := snapshotter.SaveTo("/nonexistent/dir/snap.json", s)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestCompare_NoDrift(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	baseline := snapshotter.Take("base", env)
	current := snapshotter.Take("cur", env)
	d := snapshotter.Compare(baseline, current)
	if d.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestCompare_DetectsDrift(t *testing.T) {
	baseline := snapshotter.Take("base", map[string]string{"A": "1", "B": "old", "C": "gone"})
	current := snapshotter.Take("cur", map[string]string{"A": "1", "B": "new", "D": "added"})
	d := snapshotter.Compare(baseline, current)
	if !d.HasDrift() {
		t.Fatal("expected drift")
	}
	if _, ok := d.Added["D"]; !ok {
		t.Error("expected D in Added")
	}
	if _, ok := d.Removed["C"]; !ok {
		t.Error("expected C in Removed")
	}
	if pair, ok := d.Changed["B"]; !ok || pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected B changed old->new, got %v", d.Changed["B"])
	}
}

func TestCompare_EmptySnapshots(t *testing.T) {
	baseline := snapshotter.Take("base", map[string]string{})
	current := snapshotter.Take("cur", map[string]string{})
	d := snapshotter.Compare(baseline, current)
	if d.HasDrift() {
		t.Error("expected no drift for empty snapshots")
	}
}

func init() {
	_ = os.Getenv // ensure os imported for potential future use
	_ = filepath.Join
}
