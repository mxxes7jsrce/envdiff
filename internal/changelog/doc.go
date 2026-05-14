// Package changelog compares two environment variable maps — a "before" and
// an "after" snapshot — and produces a structured Changelog describing every
// key that was added, removed, changed, or left unchanged.
//
// Typical usage:
//
//	before, _ := parser.ParseFile(".env.baseline")
//	after, _ := parser.ParseFile(".env.current")
//	cl := changelog.Build(before, after)
//	fmt.Println(cl.Summary())
//
// The Changelog type exposes a slice of Entry values sorted alphabetically
// by key, making output deterministic across runs.
package changelog
