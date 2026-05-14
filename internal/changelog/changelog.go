// Package changelog generates a human-readable changelog by comparing
// two snapshots of environment variable state over time.
package changelog

import (
	"fmt"
	"sort"
	"time"
)

// EntryKind describes the type of change recorded in a changelog entry.
type EntryKind string

const (
	Added    EntryKind = "added"
	Removed  EntryKind = "removed"
	Changed  EntryKind = "changed"
	Unchanged EntryKind = "unchanged"
)

// Entry represents a single key-level change between two env snapshots.
type Entry struct {
	Key      string
	Kind     EntryKind
	OldValue string
	NewValue string
}

// Changelog holds a timestamped set of env diff entries.
type Changelog struct {
	GeneratedAt time.Time
	Entries     []Entry
}

// Build compares two env maps (before, after) and returns a Changelog
// describing every key that was added, removed, changed, or unchanged.
func Build(before, after map[string]string) Changelog {
	seen := make(map[string]bool)
	var entries []Entry

	for k, oldVal := range before {
		seen[k] = true
		newVal, exists := after[k]
		switch {
		case !exists:
			entries = append(entries, Entry{Key: k, Kind: Removed, OldValue: oldVal})
		case oldVal != newVal:
			entries = append(entries, Entry{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		default:
			entries = append(entries, Entry{Key: k, Kind: Unchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, newVal := range after {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Kind: Added, NewValue: newVal})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Changelog{
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
	}
}

// Summary returns a short human-readable description of the changelog.
func (c Changelog) Summary() string {
	var added, removed, changed int
	for _, e := range c.Entries {
		switch e.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed", added, removed, changed)
}
