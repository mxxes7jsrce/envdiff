// Package watcher provides lightweight file-change detection for .env files.
//
// It uses MD5 checksums to determine whether file contents have changed
// between successive Poll calls, making it suitable for watch-mode tooling
// where envdiff should re-run automatically when an environment file is saved.
//
// Basic usage:
//
//	w, err := watcher.New([]string{".env", ".env.production"})
//	if err != nil { ... }
//	for _, event := range w.Poll() {
//		if event.Changed { fmt.Println(event.Path, "changed") }
//	}
package watcher
