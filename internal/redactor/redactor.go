// Package redactor provides utilities for masking sensitive values
// in .env file comparisons before display or export.
package redactor

import "strings"

// DefaultSensitivePatterns is a list of substrings that, when found in a key
// name (case-insensitive), cause the value to be redacted.
var DefaultSensitivePatterns = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private",
	"credential",
}

const redactedPlaceholder = "***REDACTED***"

// Redactor holds configuration for value masking.
type Redactor struct {
	patterns []string
}

// New returns a Redactor using the provided patterns.
// If patterns is nil or empty, DefaultSensitivePatterns is used.
func New(patterns []string) *Redactor {
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}
	return &Redactor{patterns: patterns}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// Mask returns the redacted placeholder if the key is sensitive,
// otherwise it returns the original value unchanged.
func (r *Redactor) Mask(key, value string) string {
	if r.IsSensitive(key) {
		return redactedPlaceholder
	}
	return value
}

// MaskMap returns a copy of env with sensitive values replaced.
func (r *Redactor) MaskMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = r.Mask(k, v)
	}
	return out
}
