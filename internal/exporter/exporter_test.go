package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/exporter"
)

var sampleEnv = map[string]string{
	"APP_NAME":  "envdiff",
	"DB_HOST":   "localhost",
	"DB_PASS":   "secret",
	"LOG_LEVEL": "info",
}

func TestExport_DotEnvFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv, exporter.Options{Format: exporter.FormatDotEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `APP_NAME="envdiff"`) {
		t.Errorf("expected APP_NAME line, got:\n%s", out)
	}
	if !strings.Contains(out, `DB_HOST="localhost"`) {
		t.Errorf("expected DB_HOST line, got:\n%s", out)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv, exporter.Options{Format: exporter.FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `export APP_NAME="envdiff"`) {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv, exporter.Options{Format: exporter.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key"`) || !strings.Contains(out, `"value"`) {
		t.Errorf("expected JSON fields, got:\n%s", out)
	}
}

func TestExport_PrefixFilter(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv, exporter.Options{
		Format: exporter.FormatDotEnv,
		Prefix: "DB_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("APP_NAME should be filtered out by prefix DB_")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("DB_HOST should be present")
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv, exporter.Options{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExport_EmptyMap(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, map[string]string{}, exporter.Options{Format: exporter.FormatDotEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}
