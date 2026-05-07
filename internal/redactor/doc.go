// Package redactor identifies and masks sensitive environment variable values
// before they are written to output or logs.
//
// A key is considered sensitive when its name contains one of the configured
// patterns (e.g. "password", "secret", "token"). Matching is
// case-insensitive. Sensitive values are replaced with a fixed placeholder
// string so that diffs remain useful without leaking credentials.
//
// Usage:
//
//	r := redactor.New(nil)          // uses DefaultSensitivePatterns
//	masked := r.MaskMap(envMap)     // returns a safe copy
package redactor
