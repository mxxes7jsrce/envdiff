// Package reporter formats and outputs diff results to the user.
package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable diff report to w.
func Report(w io.Writer, result diff.Result, leftName, rightName string) {
	if result.IsClean() {
		fmt.Fprintf(w, "✓ No differences found between %s and %s\n", leftName, rightName)
		return
	}

	if len(result.MissingInRight) > 0 {
		keys := sortedKeys(result.MissingInRight)
		fmt.Fprintf(w, "Keys present in %s but missing in %s:\n", leftName, rightName)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(result.MissingInLeft) > 0 {
		keys := sortedKeys(result.MissingInLeft)
		fmt.Fprintf(w, "Keys present in %s but missing in %s:\n", rightName, leftName)
		for _, k := range keys {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(result.Mismatched) > 0 {
		keys := make([]string, 0, len(result.Mismatched))
		for k := range result.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintf(w, "Keys with differing values:\n")
		for _, k := range keys {
			pair := result.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s\n", k)
			fmt.Fprintf(w, "      %s: %q\n", leftName, pair.Left)
			fmt.Fprintf(w, "      %s: %q\n", rightName, pair.Right)
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
