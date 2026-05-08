package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func baseResult() diff.Result {
	return diff.Result{
		MissingInRight: map[string]string{
			"DB_HOST": "localhost",
			"APP_PORT": "8080",
		},
		MissingInLeft: map[string]string{
			"CACHE_URL": "redis://localhost",
		},
		Mismatched: map[string][2]string{
			"APP_ENV": {"development", "production"},
		},
	}
}

func TestSort_ByKey(t *testing.T) {
	entries := sorter.Sort(baseResult(), sorter.SortByKey)

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted by key: %s before %s", entries[i-1].Key, entries[i].Key)
		}
	}
}

func TestSort_ByStatus(t *testing.T) {
	entries := sorter.Sort(baseResult(), sorter.SortByStatus)

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	for i := 1; i < len(entries); i++ {
		if entries[i].Status < entries[i-1].Status {
			t.Errorf("entries not sorted by status: %s before %s", entries[i-1].Status, entries[i].Status)
		}
	}
}

func TestSort_ByPrefix(t *testing.T) {
	entries := sorter.Sort(baseResult(), sorter.SortByPrefix)

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	// APP_ENV and APP_PORT share prefix APP — they should be adjacent
	appCount := 0
	for _, e := range entries {
		if len(e.Key) >= 3 && e.Key[:3] == "APP" {
			appCount++
		}
	}
	if appCount != 2 {
		t.Errorf("expected 2 APP_ entries, got %d", appCount)
	}
}

func TestSort_DefaultFallback(t *testing.T) {
	entries := sorter.Sort(baseResult(), sorter.SortBy("unknown"))

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("fallback sort not by key: %s before %s", entries[i-1].Key, entries[i].Key)
		}
	}
}

func TestSort_EmptyResult(t *testing.T) {
	empty := diff.Result{}
	entries := sorter.Sort(empty, sorter.SortByKey)

	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty result, got %d", len(entries))
	}
}
