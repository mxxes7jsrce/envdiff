package annotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
)

func baseResult() diff.Result {
	return diff.Result{
		MissingInRight: []string{"DB_HOST"},
		MissingInLeft:  []string{"REDIS_URL"},
		Mismatched: map[string]diff.ValuePair{
			"APP_ENV": {Left: "development", Right: "production"},
		},
	}
}

func TestAnnotate_MissingInRight(t *testing.T) {
	r := annotator.Annotate(baseResult())
	found := false
	for _, a := range r.Annotations {
		if a.Key == "DB_HOST" && a.Status == "missing_in_right" {
			found = true
			if a.Explanation == "" {
				t.Error("expected non-empty explanation")
			}
			if a.Suggestion == "" {
				t.Error("expected non-empty suggestion")
			}
		}
	}
	if !found {
		t.Error("expected annotation for DB_HOST missing_in_right")
	}
}

func TestAnnotate_MissingInLeft(t *testing.T) {
	r := annotator.Annotate(baseResult())
	found := false
	for _, a := range r.Annotations {
		if a.Key == "REDIS_URL" && a.Status == "missing_in_left" {
			found = true
		}
	}
	if !found {
		t.Error("expected annotation for REDIS_URL missing_in_left")
	}
}

func TestAnnotate_Mismatched(t *testing.T) {
	r := annotator.Annotate(baseResult())
	found := false
	for _, a := range r.Annotations {
		if a.Key == "APP_ENV" && a.Status == "mismatched" {
			found = true
			if a.Explanation == "" {
				t.Error("expected explanation for mismatched key")
			}
		}
	}
	if !found {
		t.Error("expected annotation for APP_ENV mismatched")
	}
}

func TestAnnotate_EmptyResult(t *testing.T) {
	r := annotator.Annotate(diff.Result{})
	if len(r.Annotations) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(r.Annotations))
	}
}

func TestAnnotate_TotalCount(t *testing.T) {
	r := annotator.Annotate(baseResult())
	// 1 missing_in_right + 1 missing_in_left + 1 mismatched = 3
	if len(r.Annotations) != 3 {
		t.Errorf("expected 3 annotations, got %d", len(r.Annotations))
	}
}
