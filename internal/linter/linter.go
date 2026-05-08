// Package linter checks .env file entries for common style and
// correctness issues such as invalid key names, empty values, and
// duplicate keys.
package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity indicates how serious a lint issue is.
type Severity string

const (
	Warn  Severity = "warn"
	Error Severity = "error"
)

// Issue represents a single lint finding.
type Issue struct {
	Line     int
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] line %d: %s — %s", i.Severity, i.Line, i.Key, i.Message)
}

// validKey matches POSIX-style env var names.
var validKey = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Lint analyses the provided key/value pairs (as parsed from a .env file)
// together with the original raw lines for line-number attribution.
// It returns a slice of Issues; an empty slice means no problems found.
func Lint(lines []string) []Issue {
	var issues []Issue
	seen := make(map[string]int) // key -> first line number

	for lineNum, raw := range lines {
		ln := lineNum + 1
		trimmed := strings.TrimSpace(raw)

		// skip blanks and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			issues = append(issues, Issue{Line: ln, Key: trimmed, Message: "missing '=' separator", Severity: Error})
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		value := strings.TrimSpace(trimmed[eqIdx+1:])

		if key == "" {
			issues = append(issues, Issue{Line: ln, Key: "(empty)", Message: "key must not be empty", Severity: Error})
			continue
		}

		if !validKey.MatchString(key) {
			issues = append(issues, Issue{Line: ln, Key: key, Message: "key contains invalid characters", Severity: Error})
		}

		if value == "" {
			issues = append(issues, Issue{Line: ln, Key: key, Message: "value is empty", Severity: Warn})
		}

		if first, dup := seen[key]; dup {
			issues = append(issues, Issue{Line: ln, Key: key,
				Message: fmt.Sprintf("duplicate key, first seen on line %d", first),
				Severity: Error})
		} else {
			seen[key] = ln
		}
	}

	return issues
}
