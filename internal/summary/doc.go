// Package summary computes aggregate statistics from a diff.Result.
//
// It is intended to give users a quick, high-level overview of how two
// .env files compare without requiring them to inspect every individual
// key difference.
//
// Example usage:
//
//	result := diff.Compare(left, right)
//	stats  := summary.Compute(result)
//	fmt.Println(stats) // "5 keys total: 1 missing in left, …"
package summary
