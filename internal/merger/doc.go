// Package merger provides functionality for combining two env maps into a
// single unified map.
//
// It supports three conflict-resolution strategies:
//
//   - PreferLeft  — on key collision, the value from the left (base) map wins.
//   - PreferRight — on key collision, the value from the right (override) map wins.
//   - SkipConflicts — keys with differing values are excluded from the output.
//
// Keys that exist in only one of the two maps are always included in the
// merged result regardless of the chosen strategy.
//
// Conflicts (keys present in both maps with different values) are collected and
// returned in the Result so callers can report or act on them.
package merger
