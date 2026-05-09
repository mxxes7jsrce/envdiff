package differ_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestDiff_NoChanges(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}
	deltas := differ.Diff(left, right)
	if len(deltas) != 0 {
		t.Fatalf("expected no deltas, got %d", len(deltas))
	}
}

func TestDiff_ValueChanged(t *testing.T) {
	left := map[string]string{"FOO": "old"}
	right := map[string]string{"FOO": "new"}
	deltas := differ.Diff(left, right)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Key != "FOO" || deltas[0].From != "old" || deltas[0].To != "new" {
		t.Errorf("unexpected delta: %+v", deltas[0])
	}
}

func TestDiff_CaseChangeHint(t *testing.T) {
	left := map[string]string{"MODE": "production"}
	right := map[string]string{"MODE": "PRODUCTION"}
	deltas := differ.Diff(left, right)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Hint != "case change only" {
		t.Errorf("expected 'case change only', got %q", deltas[0].Hint)
	}
}

func TestDiff_EmptyToSet(t *testing.T) {
	left := map[string]string{"TOKEN": ""}
	right := map[string]string{"TOKEN": "abc123"}
	deltas := differ.Diff(left, right)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Hint != "was empty, now set" {
		t.Errorf("unexpected hint: %q", deltas[0].Hint)
	}
}

func TestDiff_MissingKeyIgnored(t *testing.T) {
	left := map[string]string{"ONLY_LEFT": "val", "SHARED": "same"}
	right := map[string]string{"SHARED": "same"}
	deltas := differ.Diff(left, right)
	if len(deltas) != 0 {
		t.Fatalf("expected no deltas for missing key, got %d", len(deltas))
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	left := map[string]string{"Z_KEY": "a", "A_KEY": "b", "M_KEY": "c"}
	right := map[string]string{"Z_KEY": "x", "A_KEY": "y", "M_KEY": "z"}
	deltas := differ.Diff(left, right)
	if len(deltas) != 3 {
		t.Fatalf("expected 3 deltas, got %d", len(deltas))
	}
	if deltas[0].Key != "A_KEY" || deltas[1].Key != "M_KEY" || deltas[2].Key != "Z_KEY" {
		t.Errorf("deltas not sorted: %v %v %v", deltas[0].Key, deltas[1].Key, deltas[2].Key)
	}
}

func TestDelta_String(t *testing.T) {
	d := differ.Delta{Key: "FOO", From: "old", To: "new", Hint: "value changed"}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
