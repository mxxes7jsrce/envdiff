// Package baseline provides functionality for saving and loading a reference
// snapshot of a parsed .env file, allowing future comparisons against a known
// good state.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a saved baseline of key-value pairs along with metadata
// describing when and from which file the snapshot was taken.
type Snapshot struct {
	File      string            `json:"file"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// Save writes the provided env map as a baseline snapshot to the given path.
// The source file name is stored as metadata alongside a creation timestamp.
func Save(path, sourceFile string, env map[string]string) error {
	snap := Snapshot{
		File:      sourceFile,
		CreatedAt: time.Now().UTC(),
		Env:       env,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write snapshot: %w", err)
	}

	return nil
}

// Load reads a previously saved baseline snapshot from the given path.
// Returns an error if the file does not exist or cannot be decoded.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: snapshot not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: read snapshot: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: decode snapshot: %w", err)
	}

	return &snap, nil
}
