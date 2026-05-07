// Package formatter writes diff results to an io.Writer in the requested format.
package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format enumerates supported output formats.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Write serialises result to w using the given format.
func Write(w io.Writer, result diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, result)
	case FormatCSV:
		return writeCSV(w, result)
	default:
		return writeText(w, result)
	}
}

func writeText(w io.Writer, result diff.Result) error {
	keys := sortedKeys(result.MissingInRight, result.MissingInLeft, result.Mismatched)
	for _, k := range keys {
		switch {
		case contains(result.MissingInRight, k):
			fmt.Fprintf(w, "MISSING_IN_RIGHT  %s\n", k)
		case contains(result.MissingInLeft, k):
			fmt.Fprintf(w, "MISSING_IN_LEFT   %s\n", k)
		default:
			m := result.Mismatched[k]
			fmt.Fprintf(w, "MISMATCH          %s  left=%s  right=%s\n", k, m.Left, m.Right)
		}
	}
	return nil
}

type jsonEntry struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	Left   string `json:"left,omitempty"`
	Right  string `json:"right,omitempty"`
}

func writeJSON(w io.Writer, result diff.Result) error {
	var entries []jsonEntry
	for _, k := range sortedKeys(result.MissingInRight, result.MissingInLeft, result.Mismatched) {
		switch {
		case contains(result.MissingInRight, k):
			entries = append(entries, jsonEntry{Key: k, Status: "missing_in_right"})
		case contains(result.MissingInLeft, k):
			entries = append(entries, jsonEntry{Key: k, Status: "missing_in_left"})
		default:
			m := result.Mismatched[k]
			entries = append(entries, jsonEntry{Key: k, Status: "mismatch", Left: m.Left, Right: m.Right})
		}
	}
	if entries == nil {
		entries = []jsonEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func writeCSV(w io.Writer, result diff.Result) error {
	fmt.Fprintln(w, "key,status,left,right")
	for _, k := range sortedKeys(result.MissingInRight, result.MissingInLeft, result.Mismatched) {
		switch {
		case contains(result.MissingInRight, k):
			fmt.Fprintf(w, "%s,missing_in_right,,\n", csvEscape(k))
		case contains(result.MissingInLeft, k):
			fmt.Fprintf(w, "%s,missing_in_left,,\n", csvEscape(k))
		default:
			m := result.Mismatched[k]
			fmt.Fprintf(w, "%s,mismatch,%s,%s\n", csvEscape(k), csvEscape(m.Left), csvEscape(m.Right))
		}
	}
	return nil
}

func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func sortedKeys(slices ...interface{}) []string {
	seen := map[string]struct{}{}
	for _, s := range slices {
		switch v := s.(type) {
		case []string:
			for _, k := range v {
				seen[k] = struct{}{}
			}
		case map[string]diff.ValuePair:
			for k := range v {
				seen[k] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
