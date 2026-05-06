// Package diff provides functionality for comparing parsed .env file maps
// and surfacing missing or mismatched keys across environments.
package diff

// Result holds the outcome of comparing two .env maps.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey represents a key whose value differs between two env maps.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare takes two env maps (key -> value) and returns a Result describing
// all differences between them.
func Compare(left, right map[string]string) Result {
	result := Result{}

	for key, leftVal := range left {
		rightVal, ok := right[key]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, key)
			continue
		}
		if leftVal != rightVal {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        key,
				LeftValue:  leftVal,
				RightValue: rightVal,
			})
		}
	}

	for key := range right {
		if _, ok := left[key]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, key)
		}
	}

	return result
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInRight) > 0 ||
		len(r.MissingInLeft) > 0 ||
		len(r.Mismatched) > 0
}
