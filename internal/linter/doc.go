// Package linter provides static analysis for .env files.
//
// It inspects raw file lines and reports Issues for:
//
//   - Keys with invalid characters (must match [A-Za-z_][A-Za-z0-9_]*)
//   - Missing '=' separator
//   - Empty values (reported as warnings)
//   - Duplicate keys within the same file
//
// Usage:
//
//	lines := strings.Split(rawContent, "\n")
//	issues := linter.Lint(lines)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package linter
