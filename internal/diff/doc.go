// Package diff implements the core comparison logic for envdiff.
//
// It accepts two maps of environment variable key-value pairs — typically
// produced by the parser package — and returns a structured Result describing:
//
//   - Keys present in the left (base) env but missing in the right (target) env.
//   - Keys present in the right (target) env but missing in the left (base) env.
//   - Keys present in both envs but whose values differ.
//
// Example usage:
//
//	left, _ := parser.ParseFile(".env.development")
//	right, _ := parser.ParseFile(".env.production")
//	result := diff.Compare(left, right)
//	if result.HasDifferences() {
//		// report differences
//	}
package diff
