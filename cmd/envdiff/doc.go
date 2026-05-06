// Package main implements the envdiff command-line tool.
//
// envdiff compares two .env files and reports keys that are missing
// in either file or whose values differ between them.
//
// Usage:
//
//	envdiff [flags] <file1> <file2>
//
// Flags:
//
//	-quiet    suppress output; rely on exit code only
//
// Exit codes:
//
//	0   files are identical (no differences found)
//	1   an error occurred (e.g. file not found, parse error)
//	2   differences were found between the two files
package main
