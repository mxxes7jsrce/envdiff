// Package filter provides utilities for filtering diff results
// based on key patterns, prefixes, or exclusion rules.
package filter

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options holds the configuration for filtering diff results.
type Options struct {
	// Prefix restricts results to keys that start with the given prefix.
	Prefix string

	// ExcludeKeys is a set of exact key names to exclude from results.
	ExcludeKeys []string

	// OnlyMissing restricts results to keys missing in either file.
	OnlyMissing bool
}

// Apply filters a diff.Result according to the provided Options and returns
// a new diff.Result containing only the matching entries.
func Apply(result diff.Result, opts Options) diff.Result {
	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	filtered := diff.Result{
		MissingInRight: filterKeys(result.MissingInRight, opts.Prefix, excluded),
		MissingInLeft:  filterKeys(result.MissingInLeft, opts.Prefix, excluded),
		Mismatched:     filterMismatched(result.Mismatched, opts.Prefix, excluded, opts.OnlyMissing),
	}

	return filtered
}

// filterKeys returns keys from the input slice that match the prefix filter
// and are not in the excluded set.
func filterKeys(keys []string, prefix string, excluded map[string]bool) []string {
	var out []string
	for _, k := range keys {
		if excluded[k] {
			continue
		}
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		out = append(out, k)
	}
	return out
}

// filterMismatched returns mismatched entries that match the prefix filter
// and are not excluded. When onlyMissing is true, mismatched entries are
// omitted entirely.
func filterMismatched(entries []diff.Mismatch, prefix string, excluded map[string]bool, onlyMissing bool) []diff.Mismatch {
	if onlyMissing {
		return nil
	}
	var out []diff.Mismatch
	for _, m := range entries {
		if excluded[m.Key] {
			continue
		}
		if prefix != "" && !strings.HasPrefix(m.Key, prefix) {
			continue
		}
		out = append(out, m)
	}
	return out
}
