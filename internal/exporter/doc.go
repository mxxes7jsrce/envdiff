// Package exporter provides utilities for serialising a key-value environment
// map into various output formats.
//
// Supported formats:
//
//   - dotenv  – standard KEY="value" pairs (default)
//   - shell   – export KEY="value" suitable for sourcing in bash/zsh
//   - json    – JSON array of {"key": ..., "value": ...} objects
//
// An optional Prefix field on Options limits output to keys that start with
// the given string, making it easy to extract a subset of variables.
package exporter
