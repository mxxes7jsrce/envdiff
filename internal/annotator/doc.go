// Package annotator enriches diff results with human-readable explanations
// and actionable suggestions for each discrepancy found between .env files.
//
// Usage:
//
//	result := diff.Compare(left, right)
//	annotated := annotator.Annotate(result)
//	for _, a := range annotated.Annotations {
//		fmt.Printf("[%s] %s\n  → %s\n", a.Status, a.Explanation, a.Suggestion)
//	}
package annotator
