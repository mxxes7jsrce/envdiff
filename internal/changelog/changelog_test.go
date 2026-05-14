package changelog_test

import (
	"testing"

	"github.com/user/envdiff/internal/changelog"
)

func TestBuild_NoChanges(t *testing.T) {
	env := map[string]string{"KEY": "value", "PORT": "8080"}
	cl := changelog.Build(env, env)

	for _, e := range cl.Entries {
		if e.Kind != changelog.Unchanged {
			t.Errorf("expected unchanged for %q, got %q", e.Key, e.Kind)
		}
	}
}

func TestBuild_AddedKey(t *testing.T) {
	before := map[string]string{"KEY": "val"}
	after := map[string]string{"KEY": "val", "NEW": "fresh"}

	cl := changelog.Build(before, after)
	found := findEntry(cl.Entries, "NEW")
	if found == nil {
		t.Fatal("expected entry for NEW")
	}
	if found.Kind != changelog.Added {
		t.Errorf("expected added, got %q", found.Kind)
	}
	if found.NewValue != "fresh" {
		t.Errorf("unexpected NewValue: %q", found.NewValue)
	}
}

func TestBuild_RemovedKey(t *testing.T) {
	before := map[string]string{"KEY": "val", "OLD": "gone"}
	after := map[string]string{"KEY": "val"}

	cl := changelog.Build(before, after)
	found := findEntry(cl.Entries, "OLD")
	if found == nil {
		t.Fatal("expected entry for OLD")
	}
	if found.Kind != changelog.Removed {
		t.Errorf("expected removed, got %q", found.Kind)
	}
}

func TestBuild_ChangedKey(t *testing.T) {
	before := map[string]string{"PORT": "3000"}
	after := map[string]string{"PORT": "4000"}

	cl := changelog.Build(before, after)
	found := findEntry(cl.Entries, "PORT")
	if found == nil {
		t.Fatal("expected entry for PORT")
	}
	if found.Kind != changelog.Changed {
		t.Errorf("expected changed, got %q", found.Kind)
	}
	if found.OldValue != "3000" || found.NewValue != "4000" {
		t.Errorf("unexpected values: old=%q new=%q", found.OldValue, found.NewValue)
	}
}

func TestBuild_SortedEntries(t *testing.T) {
	before := map[string]string{"ZEBRA": "z", "ALPHA": "a"}
	after := map[string]string{"ZEBRA": "z", "ALPHA": "a"}

	cl := changelog.Build(before, after)
	if len(cl.Entries) >= 2 && cl.Entries[0].Key > cl.Entries[1].Key {
		t.Error("entries not sorted alphabetically")
	}
}

func TestSummary(t *testing.T) {
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "9", "C": "3"}

	cl := changelog.Build(before, after)
	s := cl.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
	// +1 added (C), -1 removed (B), ~1 changed (A)
	expected := "+1 added, -1 removed, ~1 changed"
	if s != expected {
		t.Errorf("summary mismatch: got %q, want %q", s, expected)
	}
}

func findEntry(entries []changelog.Entry, key string) *changelog.Entry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
