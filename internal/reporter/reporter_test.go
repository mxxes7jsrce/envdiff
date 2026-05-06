package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reporter"
)

func TestReport_Clean(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}
	reporter.Report(&buf, result, ".env.dev", ".env.prod")

	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestReport_MissingInRight(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInRight: map[string]string{"SECRET_KEY": "abc"},
	}
	reporter.Report(&buf, result, ".env.dev", ".env.prod")
	out := buf.String()

	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, ".env.prod") {
		t.Errorf("expected right file name in output, got: %s", out)
	}
}

func TestReport_MissingInLeft(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInLeft: map[string]string{"NEW_FEATURE": "true"},
	}
	reporter.Report(&buf, result, ".env.dev", ".env.prod")
	out := buf.String()

	if !strings.Contains(out, "NEW_FEATURE") {
		t.Errorf("expected NEW_FEATURE in output, got: %s", out)
	}
}

func TestReport_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Mismatched: map[string]diff.ValuePair{
			"DB_HOST": {Left: "localhost", Right: "db.prod.example.com"},
		},
	}
	reporter.Report(&buf, result, ".env.dev", ".env.prod")
	out := buf.String()

	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected left value in output, got: %s", out)
	}
	if !strings.Contains(out, "db.prod.example.com") {
		t.Errorf("expected right value in output, got: %s", out)
	}
}
