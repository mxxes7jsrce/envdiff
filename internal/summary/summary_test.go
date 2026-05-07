package summary_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func baseResult() diff.Result {
	return diff.Result{
		Common: map[string]string{
			"HOST": "localhost",
			"PORT": "8080",
		},
		MissingInLeft:  []string{"SECRET"},
		MissingInRight: []string{"DEBUG"},
		Mismatched: map[string][2]string{
			"DB_URL": {"postgres://a", "postgres://b"},
		},
	}
}

func TestCompute_Counts(t *testing.T) {
	s := summary.Compute(baseResult())

	if s.TotalKeys != 5 {
		t.Errorf("TotalKeys: want 5, got %d", s.TotalKeys)
	}
	if s.MissingInLeft != 1 {
		t.Errorf("MissingInLeft: want 1, got %d", s.MissingInLeft)
	}
	if s.Mismatched != 1 {
		t.Errorf("Mismatched: want 1, got %d", s.Mismatched)
	}
	if s.Clean {
		t.Error("Clean: want false, got true")
	}
}

func TestCompute_Clean(t *testing.T) {
	r := diff.Result{
		Common: map[string]string{"KEY": "val"},
	}
	s := summary.Compute(r)

	if !s.Clean {
		t.Error("Clean: want true, got false")
	}
	if s.TotalKeys != 1 {
		t.Errorf("TotalKeys: want 1, got %d", s.TotalKeys)
	}
}

func TestStats_String_Clean(t *testing.T) {
	s := summary.Stats{TotalKeys: 3, Clean: true}
	got := s.String()
	if !strings.Contains(got, "3 keys match") {
		t.Errorf("unexpected string: %q", got)
	}
}

func TestStats_String_Dirty(t *testing.T) {
	s := summary.Stats{TotalKeys: 5, MissingInLeft: 1, MissingInRight: 1, Mismatched: 1}
	got := s.String()
	if !strings.Contains(got, "5 keys total") {
		t.Errorf("unexpected string: %q", got)
	}
}
