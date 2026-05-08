// Package profiler analyzes .env files and produces a health profile
// summarizing key counts, coverage, and potential issues.
package profiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Profile holds the results of profiling two env maps against each other.
type Profile struct {
	TotalKeys      int
	PresentInBoth  int
	MissingInLeft  int
	MissingInRight int
	Mismatched     int
	CoveragePercent float64
	TopMissingKeys []string
}

// String returns a human-readable summary of the profile.
func (p Profile) String() string {
	return fmt.Sprintf(
		"total=%d both=%d missingLeft=%d missingRight=%d mismatched=%d coverage=%.1f%%",
		p.TotalKeys, p.PresentInBoth, p.MissingInLeft, p.MissingInRight,
		p.Mismatched, p.CoveragePercent,
	)
}

// Compute builds a Profile from two env maps (left and right).
func Compute(left, right map[string]string) Profile {
	result := diff.Compare(left, right)

	total := uniqueKeyCount(left, right)
	presentInBoth := total - len(result.MissingInLeft) - len(result.MissingInRight)

	var coverage float64
	if total > 0 {
		coverage = float64(presentInBoth) / float64(total) * 100.0
	}

	topMissing := collectTopMissing(result.MissingInLeft, result.MissingInRight, 5)

	return Profile{
		TotalKeys:       total,
		PresentInBoth:   presentInBoth,
		MissingInLeft:   len(result.MissingInLeft),
		MissingInRight:  len(result.MissingInRight),
		Mismatched:      len(result.Mismatched),
		CoveragePercent: coverage,
		TopMissingKeys:  topMissing,
	}
}

func uniqueKeyCount(left, right map[string]string) int {
	seen := make(map[string]struct{}, len(left)+len(right))
	for k := range left {
		seen[k] = struct{}{}
	}
	for k := range right {
		seen[k] = struct{}{}
	}
	return len(seen)
}

func collectTopMissing(missingLeft, missingRight map[string]string, n int) []string {
	keys := make([]string, 0, len(missingLeft)+len(missingRight))
	for k := range missingLeft {
		keys = append(keys, k)
	}
	for k := range missingRight {
		if !strings.Contains(strings.Join(keys, "|"), k) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	if len(keys) > n {
		return keys[:n]
	}
	return keys
}
