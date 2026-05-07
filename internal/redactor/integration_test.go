package redactor_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
	"github.com/user/envdiff/internal/redactor"
)

// TestRedactor_WithFormatter verifies that sensitive values are masked before
// being handed to the formatter, so they never appear in output.
func TestRedactor_WithFormatter(t *testing.T) {
	left := map[string]string{
		"APP_ENV":     "staging",
		"DB_PASSWORD": "old_secret",
		"PORT":        "3000",
	}
	right := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "new_secret",
		"PORT":        "3000",
	}

	r := redactor.New(nil)
	maskedLeft := r.MaskMap(left)
	maskedRight := r.MaskMap(right)

	result := diff.Compare(maskedLeft, maskedRight)

	var buf bytes.Buffer
	if err := formatter.Write(&buf, result, formatter.FormatText); err != nil {
		t.Fatalf("formatter.Write: %v", err)
	}

	output := buf.String()

	if strings.Contains(output, "old_secret") || strings.Contains(output, "new_secret") {
		t.Errorf("raw secrets must not appear in output; got:\n%s", output)
	}
	if !strings.Contains(output, "APP_ENV") {
		t.Errorf("expected APP_ENV mismatch in output; got:\n%s", output)
	}
}
