// Package sorter provides utilities for sorting and grouping diff results
// by key name, status, or prefix for structured output and reporting.
package sorter

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// SortBy defines the field used to order diff entries.
type SortBy string

const (
	SortByKey    SortBy = "key"
	SortByStatus SortBy = "status"
	SortByPrefix SortBy = "prefix"
)

// Entry represents a single diff item with its status and key.
type Entry struct {
	Key    string
	Status string
	Left   string
	Right  string
}

// Sort takes a diff.Result and returns a slice of Entry values ordered
// according to the given SortBy field. Unknown values fall back to SortByKey.
func Sort(result diff.Result, by SortBy) []Entry {
	entries := flatten(result)

	switch by {
	case SortByStatus:
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].Status != entries[j].Status {
				return entries[i].Status < entries[j].Status
			}
			return entries[i].Key < entries[j].Key
		})
	case SortByPrefix:
		sort.SliceStable(entries, func(i, j int) bool {
			pi := prefix(entries[i].Key)
			pj := prefix(entries[j].Key)
			if pi != pj {
				return pi < pj
			}
			return entries[i].Key < entries[j].Key
		})
	default:
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
	}

	return entries
}

// flatten converts a diff.Result into a flat list of Entry values.
func flatten(result diff.Result) []Entry {
	var entries []Entry

	for k, v := range result.MissingInRight {
		entries = append(entries, Entry{Key: k, Status: "missing_right", Left: v, Right: ""})
	}
	for k, v := range result.MissingInLeft {
		entries = append(entries, Entry{Key: k, Status: "missing_left", Left: "", Right: v})
	}
	for k, p := range result.Mismatched {
		entries = append(entries, Entry{Key: k, Status: "mismatched", Left: p[0], Right: p[1]})
	}

	return entries
}

// prefix returns the portion of a key before the first underscore,
// or the full key if no underscore is present.
func prefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
