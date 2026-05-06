// Package reporter provides functionality for formatting and displaying
// the results of an environment file comparison.
//
// It accepts a diff.Result and renders it to any io.Writer in a
// human-readable text format, clearly indicating:
//
//   - Keys missing from the right-hand file
//   - Keys missing from the left-hand file
//   - Keys present in both files but with differing values
//
// Example usage:
//
//	reporter.Report(os.Stdout, result, ".env.development", ".env.production")
package reporter
