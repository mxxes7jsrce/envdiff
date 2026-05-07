package merger

import "sort"

// Strategy controls how conflicting keys are resolved during a merge.
type Strategy int

const (
	// PreferLeft keeps the value from the left map on conflict.
	PreferLeft Strategy = iota
	// PreferRight keeps the value from the right map on conflict.
	PreferRight
	// SkipConflicts omits any key that exists in both maps with different values.
	SkipConflicts
)

// Result holds the merged key-value pairs and metadata about the operation.
type Result struct {
	Merged    map[string]string
	Conflicts []string // keys that had differing values
}

// Merge combines left and right env maps according to the given strategy.
// Keys present in only one map are always included unchanged.
func Merge(left, right map[string]string, strategy Strategy) Result {
	merged := make(map[string]string)
	var conflicts []string

	// Copy all keys from left.
	for k, v := range left {
		merged[k] = v
	}

	for k, rv := range right {
		lv, exists := merged[k]
		if !exists {
			// Key only in right — always include.
			merged[k] = rv
			continue
		}
		if lv == rv {
			// Identical values — no conflict.
			continue
		}
		// Conflict: values differ.
		conflicts = append(conflicts, k)
		switch strategy {
		case PreferRight:
			merged[k] = rv
		case SkipConflicts:
			delete(merged, k)
		default: // PreferLeft — keep existing value.
		}
	}

	sort.Strings(conflicts)
	return Result{Merged: merged, Conflicts: conflicts}
}
