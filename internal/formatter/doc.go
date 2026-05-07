// Package formatter provides output formatting for envdiff comparison results.
//
// It supports multiple output formats:
//
//   - text: human-readable plain-text output (default)
//   - json: structured JSON suitable for machine consumption
//   - csv:  comma-separated values for spreadsheet or pipeline use
//
// Usage:
//
//	result := diff.Compare(left, right)
//	formatter.Write(os.Stdout, result, formatter.FormatJSON)
package formatter
