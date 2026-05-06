package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

func baseResult() diff.Result {
	return diff.Result{
		MissingInRight: []string{"DB_HOST", "APP_PORT"},
		MissingInLeft:  []string{"REDIS_URL"},
		Mismatched: []diff.Mismatch{
			{Key: "APP_ENV", LeftVal: "dev", RightVal: "prod"},
			{Key: "DB_NAME", LeftVal: "mydb", RightVal: "proddb"},
		},
	}
}

func TestApply_NoOptions(t *testing.T) {
	result := filter.Apply(baseResult(), filter.Options{})
	if len(result.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(result.MissingInRight))
	}
	if len(result.MissingInLeft) != 1 {
		t.Errorf("expected 1 MissingInLeft, got %d", len(result.MissingInLeft))
	}
	if len(result.Mismatched) != 2 {
		t.Errorf("expected 2 Mismatched, got %d", len(result.Mismatched))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	result := filter.Apply(baseResult(), filter.Options{Prefix: "DB_"})
	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "DB_HOST" {
		t.Errorf("expected only DB_HOST in MissingInRight, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no MissingInLeft with DB_ prefix, got %v", result.MissingInLeft)
	}
	if len(result.Mismatched) != 1 || result.Mismatched[0].Key != "DB_NAME" {
		t.Errorf("expected only DB_NAME in Mismatched, got %v", result.Mismatched)
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	opts := filter.Options{ExcludeKeys: []string{"APP_PORT", "APP_ENV"}}
	result := filter.Apply(baseResult(), opts)
	for _, k := range result.MissingInRight {
		if k == "APP_PORT" {
			t.Error("APP_PORT should have been excluded from MissingInRight")
		}
	}
	for _, m := range result.Mismatched {
		if m.Key == "APP_ENV" {
			t.Error("APP_ENV should have been excluded from Mismatched")
		}
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	result := filter.Apply(baseResult(), filter.Options{OnlyMissing: true})
	if len(result.Mismatched) != 0 {
		t.Errorf("expected no Mismatched entries with OnlyMissing, got %d", len(result.Mismatched))
	}
	if len(result.MissingInRight) != 2 {
		t.Errorf("expected MissingInRight to be unaffected, got %d", len(result.MissingInRight))
	}
}

func TestApply_EmptyResult(t *testing.T) {
	result := filter.Apply(diff.Result{}, filter.Options{Prefix: "APP_"})
	if len(result.MissingInRight) != 0 || len(result.MissingInLeft) != 0 || len(result.Mismatched) != 0 {
		t.Error("expected all empty slices for empty input")
	}
}
