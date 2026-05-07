// Package validator provides rule-based validation for environment variable maps.
//
// Rules can enforce that certain keys are present (Required) and that their
// values match a given regular expression pattern (Pattern). Violations are
// collected and returned so callers can decide how to surface them.
//
// Example usage:
//
//	rules := []validator.Rule{
//		{Key: "APP_ENV", Required: true, Pattern: `^(production|staging)$`},
//		{Key: "PORT",    Required: true, Pattern: `^\d+$`},
//	}
//	violations := validator.Validate("prod.env", envMap, rules)
package validator
