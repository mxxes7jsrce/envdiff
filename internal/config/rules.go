package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/validator"
)

// rulesFile is the on-disk representation of validation rules.
type rulesFile struct {
	Rules []ruleEntry `json:"rules"`
}

type ruleEntry struct {
	Key      string `json:"key"`
	Pattern  string `json:"pattern,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// LoadRules reads a JSON rules file from path and returns a slice of
// validator.Rule. Returns an empty slice (not an error) when path is empty.
func LoadRules(path string) ([]validator.Rule, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading rules file %q: %w", path, err)
	}

	var rf rulesFile
	if err := json.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("parsing rules file %q: %w", path, err)
	}

	rules := make([]validator.Rule, 0, len(rf.Rules))
	for _, e := range rf.Rules {
		if e.Key == "" {
			continue
		}
		rules = append(rules, validator.Rule{
			Key:      e.Key,
			Pattern:  e.Pattern,
			Required: e.Required,
		})
	}
	return rules, nil
}
