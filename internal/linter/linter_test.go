package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestLint_NoIssues(t *testing.T) {
	lines := []string{
		"# comment",
		"",
		"APP_ENV=production",
		"PORT=8080",
	}
	issues := linter.Lint(lines)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	lines := []string{"SECRET_KEY="}
	issues := linter.Lint(lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != linter.Warn {
		t.Errorf("expected Warn severity, got %s", issues[0].Severity)
	}
	if issues[0].Key != "SECRET_KEY" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
}

func TestLint_InvalidKeyName(t *testing.T) {
	lines := []string{"123BAD=value"}
	issues := linter.Lint(lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != linter.Error {
		t.Errorf("expected Error severity, got %s", issues[0].Severity)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	lines := []string{
		"DB_HOST=localhost",
		"DB_PORT=5432",
		"DB_HOST=remotehost",
	}
	issues := linter.Lint(lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Line != 3 {
		t.Errorf("expected issue on line 3, got %d", issues[0].Line)
	}
	if issues[0].Severity != linter.Error {
		t.Errorf("expected Error severity, got %s", issues[0].Severity)
	}
}

func TestLint_MissingSeparator(t *testing.T) {
	lines := []string{"BADLINE"}
	issues := linter.Lint(lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != linter.Error {
		t.Errorf("expected Error severity, got %s", issues[0].Severity)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	lines := []string{
		"GOOD=value",
		"123BAD=x",
		"EMPTY_VAL=",
		"GOOD=again",
	}
	issues := linter.Lint(lines)
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d: %v", len(issues), issues)
	}
}

func TestIssue_String(t *testing.T) {
	i := linter.Issue{Line: 5, Key: "FOO", Message: "value is empty", Severity: linter.Warn}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
