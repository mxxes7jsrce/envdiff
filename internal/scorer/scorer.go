// Package scorer provides a similarity scoring mechanism for comparing
// two env maps and producing a numeric health score (0–100).
package scorer

import (
	"github.com/user/envdiff/internal/diff"
)

// Score holds the numeric result and a human-readable grade.
type Score struct {
	// Value is the computed score in the range [0, 100].
	Value int

	// Grade is a letter grade derived from Value.
	Grade string

	// Total is the total number of unique keys across both maps.
	Total int

	// Matched is the number of keys that are present and identical in both maps.
	Matched int
}

// String returns a short human-readable representation of the score.
func (s Score) String() string {
	return s.Grade
}

// Compute derives a Score by comparing two env maps.
// A key contributes its full weight only when it is present in both maps
// with identical values. Missing or mismatched keys reduce the score.
func Compute(left, right map[string]string) Score {
	result := diff.Compare(left, right)

	total := len(result.MissingInLeft) + len(result.MissingInRight) + len(result.Mismatched)
	// Count keys that are present and matching in both maps.
	matched := 0
	for k, lv := range left {
		if rv, ok := right[k]; ok && lv == rv {
			matched++
			total++
		}
	}

	if total == 0 {
		return Score{Value: 100, Grade: "A", Total: 0, Matched: 0}
	}

	value := int(float64(matched) / float64(total) * 100)

	return Score{
		Value:   value,
		Grade:   grade(value),
		Total:   total,
		Matched: matched,
	}
}

// grade converts a numeric score into a letter grade.
func grade(v int) string {
	switch {
	case v >= 90:
		return "A"
	case v >= 75:
		return "B"
	case v >= 60:
		return "C"
	case v >= 40:
		return "D"
	default:
		return "F"
	}
}
