package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule applied to env key-value pairs.
type Rule struct {
	Key     string
	Pattern string
	Required bool
}

// Violation describes a single validation failure.
type Violation struct {
	Key     string
	File    string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.File, v.Key, v.Message)
}

// Validate checks the provided env map against a set of rules.
// It returns a slice of Violations for any rules that are not satisfied.
func Validate(filename string, env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				File:    filename,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, val)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					File:    filename,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !matched {
				violations = append(violations, Violation{
					Key:     rule.Key,
					File:    filename,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}

		if strings.TrimSpace(val) == "" && rule.Required {
			violations = append(violations, Violation{
				Key:     rule.Key,
				File:    filename,
				Message: "required key has empty value",
			})
		}
	}

	return violations
}
