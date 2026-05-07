package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for diff results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// jsonResult is the JSON-serializable representation of a diff result.
type jsonResult struct {
	MissingInLeft  []string            `json:"missing_in_left,omitempty"`
	MissingInRight []string            `json:"missing_in_right,omitempty"`
	Mismatched     map[string][2]string `json:"mismatched,omitempty"`
}

// Write formats the diff result in the requested format and writes it to w.
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
	for _, k := range result.MissingInLeft {
		if _, err := fmt.Fprintf(w, "< missing: %s\n", k); err != nil {
			return err
		}
	}
	for _, k := range result.MissingInRight {
		if _, err := fmt.Fprintf(w, "> missing: %s\n", k); err != nil {
			return err
		}
	}
	for k, v := range result.Mismatched {
		if _, err := fmt.Fprintf(w, "~ mismatch: %s (%q vs %q)\n", k, v[0], v[1]); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, result diff.Result) error {
	out := jsonResult{
		MissingInLeft:  result.MissingInLeft,
		MissingInRight: result.MissingInRight,
	}
	if len(result.Mismatched) > 0 {
		out.Mismatched = make(map[string][2]string, len(result.Mismatched))
		for k, v := range result.Mismatched {
			out.Mismatched[k] = v
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func writeCSV(w io.Writer, result diff.Result) error {
	if _, err := fmt.Fprintln(w, "type,key,left_value,right_value"); err != nil {
		return err
	}
	for _, k := range result.MissingInLeft {
		if _, err := fmt.Fprintf(w, "missing_in_left,%s,,\n", csvEscape(k)); err != nil {
			return err
		}
	}
	for _, k := range result.MissingInRight {
		if _, err := fmt.Fprintf(w, "missing_in_right,%s,,\n", csvEscape(k)); err != nil {
			return err
		}
	}
	for k, v := range result.Mismatched {
		if _, err := fmt.Fprintf(w, "mismatch,%s,%s,%s\n", csvEscape(k), csvEscape(v[0]), csvEscape(v[1])); err != nil {
			return err
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
