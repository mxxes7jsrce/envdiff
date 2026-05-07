package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestValidate_NoViolations(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
	}
	rules := []validator.Rule{
		{Key: "APP_ENV", Required: true, Pattern: `^(production|staging|development)$`},
		{Key: "LOG_LEVEL", Required: true},
	}
	violations := validator.Validate("prod.env", env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	env := map[string]string{}
	rules := []validator.Rule{
		{Key: "DATABASE_URL", Required: true},
	}
	violations := validator.Validate("prod.env", env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key in violation: %s", violations[0].Key)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{
		"PORT": "not-a-number",
	}
	rules := []validator.Rule{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	violations := validator.Validate("dev.env", env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_OptionalKeyAbsent(t *testing.T) {
	env := map[string]string{}
	rules := []validator.Rule{
		{Key: "OPTIONAL_FEATURE", Required: false, Pattern: `^(true|false)$`},
	}
	violations := validator.Validate("dev.env", env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for absent optional key, got %d", len(violations))
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "unknown",
	}
	rules := []validator.Rule{
		{Key: "APP_ENV", Required: true, Pattern: `^(production|staging)$`},
		{Key: "SECRET_KEY", Required: true},
	}
	violations := validator.Validate("staging.env", env, rules)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(violations))
	}
}

func TestViolation_String(t *testing.T) {
	v := validator.Violation{Key: "FOO", File: "test.env", Message: "required key is missing"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string from Violation.String()")
	}
}
