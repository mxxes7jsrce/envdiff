package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
)

func baseResult() diff.Result {
	return diff.Result{
		MissingInLeft:  []string{"DB_HOST"},
		MissingInRight: []string{"API_KEY"},
		Mismatched:     map[string][2]string{"PORT": {"8080", "9090"}},
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, baseResult(), formatter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "< missing: DB_HOST") {
		t.Errorf("expected missing-in-left entry, got:\n%s", out)
	}
	if !strings.Contains(out, "> missing: API_KEY") {
		t.Errorf("expected missing-in-right entry, got:\n%s", out)
	}
	if !strings.Contains(out, "~ mismatch: PORT") {
		t.Errorf("expected mismatch entry, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, baseResult(), formatter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"missing_in_left"`) {
		t.Errorf("expected JSON key missing_in_left, got:\n%s", out)
	}
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Errorf("expected DB_HOST in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"mismatched"`) {
		t.Errorf("expected mismatched in JSON output, got:\n%s", out)
	}
}

func TestWrite_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, baseResult(), formatter.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "type,key,left_value,right_value" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if len(lines) < 4 {
		t.Errorf("expected at least 4 CSV lines, got %d", len(lines))
	}
}

func TestWrite_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	empty := diff.Result{}
	for _, f := range []formatter.Format{formatter.FormatText, formatter.FormatJSON, formatter.FormatCSV} {
		buf.Reset()
		if err := formatter.Write(&buf, empty, f); err != nil {
			t.Errorf("format %s: unexpected error: %v", f, err)
		}
	}
}

func TestWrite_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	// unknown format should fall back to text
	if err := formatter.Write(&buf, baseResult(), formatter.Format("unknown")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output for default format fallback")
	}
}
