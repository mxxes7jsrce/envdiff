package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/scorer"
)

func TestCompute_PerfectMatch(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3"}
	right := map[string]string{"A": "1", "B": "2", "C": "3"}

	s := scorer.Compute(left, right)

	if s.Value != 100 {
		t.Errorf("expected 100, got %d", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected grade A, got %s", s.Grade)
	}
	if s.Matched != 3 {
		t.Errorf("expected 3 matched, got %d", s.Matched)
	}
}

func TestCompute_EmptyMaps(t *testing.T) {
	s := scorer.Compute(map[string]string{}, map[string]string{})

	if s.Value != 100 {
		t.Errorf("expected 100 for empty maps, got %d", s.Value)
	}
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestCompute_AllMissing(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"C": "3", "D": "4"}

	s := scorer.Compute(left, right)

	if s.Value != 0 {
		t.Errorf("expected 0, got %d", s.Value)
	}
	if s.Grade != "F" {
		t.Errorf("expected grade F, got %s", s.Grade)
	}
	if s.Matched != 0 {
		t.Errorf("expected 0 matched, got %d", s.Matched)
	}
}

func TestCompute_PartialMatch(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"}
	right := map[string]string{"A": "1", "B": "2", "C": "changed"}

	s := scorer.Compute(left, right)

	// matched=2 (A,B), total=2+1(mismatch C)+1(missing D in right)=4 → 50
	if s.Value != 50 {
		t.Errorf("expected 50, got %d", s.Value)
	}
	if s.Grade != "D" {
		t.Errorf("expected grade D, got %s", s.Grade)
	}
}

func TestScore_String(t *testing.T) {
	s := scorer.Score{Grade: "B"}
	if s.String() != "B" {
		t.Errorf("expected B, got %s", s.String())
	}
}

func TestCompute_GradeBoundaries(t *testing.T) {
	tests := []struct {
		name          string
		matched, total int
		wantGrade     string
	}{
		{"A boundary", 9, 10, "A"},  // 90
		{"B boundary", 8, 10, "B"},  // 80
		{"C boundary", 7, 10, "C"},  // 70
		{"D boundary", 6, 10, "D"},  // 60
		{"F boundary", 5, 10, "F"},  // 50
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := make(map[string]string, tt.total)
			right := make(map[string]string, tt.total)
			for i := 0; i < tt.matched; i++ {
				key := string(rune('A' + i))
				left[key] = "v"
				right[key] = "v"
			}
			for i := tt.matched; i < tt.total; i++ {
				key := string(rune('A' + i))
				left[key] = "v"
				right[key] = "different"
			}
			s := scorer.Compute(left, right)
			if s.Grade != tt.wantGrade {
				t.Errorf("expected grade %s, got %s (value=%d)", tt.wantGrade, s.Grade, s.Value)
			}
		})
	}
}
