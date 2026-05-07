package watcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// FileState holds the last known checksum and modification time of a file.
type FileState struct {
	Path    string
	Checksum string
	ModTime  time.Time
}

// ChangeEvent describes a detected change between two Watch calls.
type ChangeEvent struct {
	Path    string
	Changed bool
	Err     error
}

// Watcher tracks one or more .env files for changes.
type Watcher struct {
	states map[string]FileState
}

// New creates a Watcher initialised with the current state of each path.
func New(paths []string) (*Watcher, error) {
	w := &Watcher{states: make(map[string]FileState, len(paths))}
	for _, p := range paths {
		state, err := snapshot(p)
		if err != nil {
			return nil, fmt.Errorf("watcher: initial snapshot of %q: %w", p, err)
		}
		w.states[p] = state
	}
	return w, nil
}

// Poll checks all tracked files and returns a ChangeEvent for each one.
func (w *Watcher) Poll() []ChangeEvent {
	events := make([]ChangeEvent, 0, len(w.states))
	for path, prev := range w.states {
		curr, err := snapshot(path)
		if err != nil {
			events = append(events, ChangeEvent{Path: path, Err: err})
			continue
		}
		changed := curr.Checksum != prev.Checksum
		if changed {
			w.states[path] = curr
		}
		events = append(events, ChangeEvent{Path: path, Changed: changed})
	}
	return events
}

func snapshot(path string) (FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FileState{}, err
	}

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileState{}, err
	}

	return FileState{
		Path:     path,
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
		ModTime:  info.ModTime(),
	}, nil
}
