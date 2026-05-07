package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func TestMerge_NoConflicts(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"C": "3", "D": "4"}

	res := merger.Merge(left, right, merger.PreferLeft)

	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", res.Conflicts)
	}
	if res.Merged["A"] != "1" || res.Merged["B"] != "2" || res.Merged["C"] != "3" || res.Merged["D"] != "4" {
		t.Errorf("unexpected merged map: %v", res.Merged)
	}
}

func TestMerge_PreferLeft(t *testing.T) {
	left := map[string]string{"A": "left"}
	right := map[string]string{"A": "right"}

	res := merger.Merge(left, right, merger.PreferLeft)

	if res.Merged["A"] != "left" {
		t.Errorf("expected 'left', got %q", res.Merged["A"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "A" {
		t.Errorf("expected conflict on A, got %v", res.Conflicts)
	}
}

func TestMerge_PreferRight(t *testing.T) {
	left := map[string]string{"X": "old"}
	right := map[string]string{"X": "new"}

	res := merger.Merge(left, right, merger.PreferRight)

	if res.Merged["X"] != "new" {
		t.Errorf("expected 'new', got %q", res.Merged["X"])
	}
}

func TestMerge_SkipConflicts(t *testing.T) {
	left := map[string]string{"KEY": "a", "SAFE": "same"}
	right := map[string]string{"KEY": "b", "SAFE": "same", "EXTRA": "x"}

	res := merger.Merge(left, right, merger.SkipConflicts)

	if _, ok := res.Merged["KEY"]; ok {
		t.Error("expected KEY to be omitted due to conflict")
	}
	if res.Merged["SAFE"] != "same" {
		t.Errorf("expected SAFE to be preserved, got %q", res.Merged["SAFE"])
	}
	if res.Merged["EXTRA"] != "x" {
		t.Errorf("expected EXTRA to be included, got %q", res.Merged["EXTRA"])
	}
}

func TestMerge_IdenticalValues_NoConflict(t *testing.T) {
	left := map[string]string{"PORT": "8080"}
	right := map[string]string{"PORT": "8080"}

	res := merger.Merge(left, right, merger.PreferLeft)

	if len(res.Conflicts) != 0 {
		t.Errorf("identical values should not be reported as conflicts")
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	res := merger.Merge(map[string]string{}, map[string]string{}, merger.PreferLeft)

	if len(res.Merged) != 0 {
		t.Errorf("expected empty merged map")
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}
