package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoDifferences(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(left, right)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"FOO": "bar", "ONLY_LEFT": "val"}
	right := map[string]string{"FOO": "bar"}

	result := diff.Compare(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no keys missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "bar", "ONLY_RIGHT": "val"}

	result := diff.Compare(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("expected ONLY_RIGHT missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"FOO": "original", "BAR": "same"}
	right := map[string]string{"FOO": "changed", "BAR": "same"}

	result := diff.Compare(left, right)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatched key, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "FOO" || m.LeftValue != "original" || m.RightValue != "changed" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := diff.Compare(map[string]string{}, map[string]string{})
	if result.HasDifferences() {
		t.Error("expected no differences for two empty maps")
	}
}
