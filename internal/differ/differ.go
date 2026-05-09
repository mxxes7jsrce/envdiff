// Package differ provides line-level diffing of individual env values,
// highlighting what changed between two versions of the same key.
package differ

import (
	"fmt"
	"strings"
)

// Delta represents a character-level or word-level difference for a single key.
type Delta struct {
	Key  string
	From string
	To   string
	Hint string // human-readable description of the change
}

// String returns a human-readable representation of the delta.
func (d Delta) String() string {
	return fmt.Sprintf("%s: %q → %q (%s)", d.Key, d.From, d.To, d.Hint)
}

// Diff computes Deltas for keys that appear in both left and right but with
// differing values. Keys absent from either side are ignored (handled by
// internal/diff).
func Diff(left, right map[string]string) []Delta {
	var deltas []Delta
	for k, lv := range left {
		rv, ok := right[k]
		if !ok || lv == rv {
			continue
		}
		deltas = append(deltas, Delta{
			Key:  k,
			From: lv,
			To:   rv,
			Hint: describeChange(lv, rv),
		})
	}
	sortDeltas(deltas)
	return deltas
}

// describeChange returns a short hint describing the nature of the value change.
func describeChange(from, to string) string {
	switch {
	case from == "" && to != "":
		return "was empty, now set"
	case from != "" && to == "":
		return "was set, now empty"
	case strings.EqualFold(from, to):
		return "case change only"
	case len(from) != len(to):
		return fmt.Sprintf("length changed %d→%d", len(from), len(to))
	default:
		return "value changed"
	}
}

// sortDeltas sorts deltas alphabetically by key for deterministic output.
func sortDeltas(deltas []Delta) {
	for i := 1; i < len(deltas); i++ {
		for j := i; j > 0 && deltas[j].Key < deltas[j-1].Key; j-- {
			deltas[j], deltas[j-1] = deltas[j-1], deltas[j]
		}
	}
}
