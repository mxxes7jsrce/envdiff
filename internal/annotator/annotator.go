// Package annotator attaches human-readable annotations to diff results,
// explaining why each entry is flagged and suggesting remediation steps.
package annotator

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Annotation holds an explanation and optional suggestion for a diff entry.
type Annotation struct {
	Key        string
	Status     string
	Explanation string
	Suggestion  string
}

// Result wraps a diff.Result with per-entry annotations.
type Result struct {
	Annotations []Annotation
}

// Annotate inspects a diff.Result and produces an annotated Result.
func Annotate(r diff.Result) Result {
	var annotations []Annotation

	for _, k := range r.MissingInRight {
		annotations = append(annotations, Annotation{
			Key:         k,
			Status:      "missing_in_right",
			Explanation: fmt.Sprintf("Key %q is defined in the left file but absent in the right file.", k),
			Suggestion:  fmt.Sprintf("Add %q to the right .env file.", k),
		})
	}

	for _, k := range r.MissingInLeft {
		annotations = append(annotations, Annotation{
			Key:         k,
			Status:      "missing_in_left",
			Explanation: fmt.Sprintf("Key %q is defined in the right file but absent in the left file.", k),
			Suggestion:  fmt.Sprintf("Add %q to the left .env file or remove it from the right.", k),
		})
	}

	for k, pair := range r.Mismatched {
		annotations = append(annotations, Annotation{
			Key:         k,
			Status:      "mismatched",
			Explanation: fmt.Sprintf("Key %q has different values: left=%q, right=%q.", k, pair.Left, pair.Right),
			Suggestion:  fmt.Sprintf("Reconcile the value of %q between environments.", k),
		})
	}

	return Result{Annotations: annotations}
}
