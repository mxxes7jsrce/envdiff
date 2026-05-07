// Package summary provides aggregation and statistics over diff results.
package summary

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Stats holds aggregate counts derived from a diff result.
type Stats struct {
	TotalKeys      int
	MissingInLeft  int
	MissingInRight int
	Mismatched     int
	Clean          bool
}

// String returns a human-readable one-line summary of the stats.
func (s Stats) String() string {
	if s.Clean {
		return fmt.Sprintf("all %d keys match", s.TotalKeys)
	}
	return fmt.Sprintf(
		"%d keys total: %d missing in left, %d missing in right, %d mismatched",
		s.TotalKeys, s.MissingInLeft, s.MissingInRight, s.Mismatched,
	)
}

// Compute derives a Stats value from a diff.Result.
func Compute(r diff.Result) Stats {
	total := uniqueKeys(r)
	missingLeft := len(r.MissingInLeft)
	missingRight := len(r.MissingInRight)
	mismatched := len(r.Mismatched)

	return Stats{
		TotalKeys:      total,
		MissingInLeft:  missingLeft,
		MissingInRight: mismatched,
		Mismatched:     mismatched,
		Clean:          missingLeft == 0 && missingRight == 0 && mismatched == 0,
	}
}

// uniqueKeys counts the total number of distinct keys across all diff buckets.
func uniqueKeys(r diff.Result) int {
	seen := make(map[string]struct{})
	for k := range r.Common {
		seen[k] = struct{}{}
	}
	for _, k := range r.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, k := range r.MissingInRight {
		seen[k] = struct{}{}
	}
	for k := range r.Mismatched {
		seen[k] = struct{}{}
	}
	return len(seen)
}
