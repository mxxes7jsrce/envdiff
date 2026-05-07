package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/redactor"
)

func TestIsSensitive_MatchesDefaults(t *testing.T) {
	r := redactor.New(nil)
	sensitive := []string{"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN", "STRIPE_API_KEY", "PRIVATE_KEY"}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	r := redactor.New(nil)
	safe := []string{"APP_ENV", "PORT", "LOG_LEVEL", "BASE_URL"}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMask_RedactsSensitive(t *testing.T) {
	r := redactor.New(nil)
	got := r.Mask("DB_PASSWORD", "supersecret")
	if got == "supersecret" {
		t.Error("expected value to be redacted")
	}
	if got == "" {
		t.Error("expected non-empty redacted placeholder")
	}
}

func TestMask_PreservesSafeValue(t *testing.T) {
	r := redactor.New(nil)
	got := r.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected %q, got %q", "production", got)
	}
}

func TestMaskMap_RedactsCorrectKeys(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "hunter2",
		"PORT":        "8080",
		"AUTH_TOKEN":  "tok_abc123",
	}
	masked := r.MaskMap(env)
	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged")
	}
	if masked["DB_PASSWORD"] == "hunter2" {
		t.Errorf("DB_PASSWORD should be redacted")
	}
	if masked["AUTH_TOKEN"] == "tok_abc123" {
		t.Errorf("AUTH_TOKEN should be redacted")
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	r := redactor.New([]string{"internal"})
	if !r.IsSensitive("INTERNAL_KEY") {
		t.Error("expected INTERNAL_KEY to be sensitive with custom pattern")
	}
	if r.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to NOT be sensitive with custom-only patterns")
	}
}
