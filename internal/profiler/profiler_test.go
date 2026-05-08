package profiler_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/profiler"
)

func TestCompute_PerfectMatch(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}

	p := profiler.Compute(left, right)

	if p.TotalKeys != 2 {
		t.Errorf("expected TotalKeys=2, got %d", p.TotalKeys)
	}
	if p.PresentInBoth != 2 {
		t.Errorf("expected PresentInBoth=2, got %d", p.PresentInBoth)
	}
	if p.CoveragePercent != 100.0 {
		t.Errorf("expected coverage=100.0, got %.1f", p.CoveragePercent)
	}
	if p.MissingInLeft != 0 || p.MissingInRight != 0 {
		t.Errorf("expected no missing keys")
	}
}

func TestCompute_MissingInRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3"}
	right := map[string]string{"A": "1"}

	p := profiler.Compute(left, right)

	if p.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", p.TotalKeys)
	}
	if p.MissingInRight != 2 {
		t.Errorf("expected MissingInRight=2, got %d", p.MissingInRight)
	}
	if p.PresentInBoth != 1 {
		t.Errorf("expected PresentInBoth=1, got %d", p.PresentInBoth)
	}
}

func TestCompute_EmptyMaps(t *testing.T) {
	p := profiler.Compute(map[string]string{}, map[string]string{})

	if p.TotalKeys != 0 {
		t.Errorf("expected TotalKeys=0, got %d", p.TotalKeys)
	}
	if p.CoveragePercent != 0.0 {
		t.Errorf("expected coverage=0.0, got %.1f", p.CoveragePercent)
	}
}

func TestCompute_TopMissingKeys_Capped(t *testing.T) {
	left := map[string]string{
		"KEY1": "a", "KEY2": "b", "KEY3": "c",
		"KEY4": "d", "KEY5": "e", "KEY6": "f",
	}
	right := map[string]string{}

	p := profiler.Compute(left, right)

	if len(p.TopMissingKeys) > 5 {
		t.Errorf("expected at most 5 top missing keys, got %d", len(p.TopMissingKeys))
	}
}

func TestProfile_String(t *testing.T) {
	left := map[string]string{"X": "1", "Y": "2"}
	right := map[string]string{"X": "1"}

	p := profiler.Compute(left, right)
	s := p.String()

	if !strings.Contains(s, "total=") || !strings.Contains(s, "coverage=") {
		t.Errorf("unexpected String() output: %s", s)
	}
}
