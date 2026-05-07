package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/config"
)

func writeTempRules(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal rules: %v", err)
	}
	p := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatalf("write rules: %v", err)
	}
	return p
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := config.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestLoadRules_ValidFile(t *testing.T) {
	payload := map[string]any{
		"rules": []map[string]any{
			{"key": "APP_ENV", "required": true, "pattern": `^(production|staging)$`},
			{"key": "PORT", "required": true, "pattern": `^\d+$`},
		},
	}
	path := writeTempRules(t, payload)
	rules, err := config.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Key != "APP_ENV" || !rules[0].Required {
		t.Errorf("unexpected first rule: %+v", rules[0])
	}
}

func TestLoadRules_SkipsEmptyKeys(t *testing.T) {
	payload := map[string]any{
		"rules": []map[string]any{
			{"key": "", "required": true},
			{"key": "VALID_KEY"},
		},
	}
	path := writeTempRules(t, payload)
	rules, err := config.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule after skipping empty key, got %d", len(rules))
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := config.LoadRules("/nonexistent/path/rules.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadRules_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(p, []byte("not json{"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	_, err := config.LoadRules(p)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
