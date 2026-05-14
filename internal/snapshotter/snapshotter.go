// Package snapshotter provides functionality for capturing and comparing
// point-in-time snapshots of parsed .env file state, enabling drift detection
// between a known-good baseline and the current environment.
package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of a parsed .env map.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Source     string            `json:"source"`
	Env        map[string]string `json:"env"`
}

// Drift describes the difference between two snapshots.
type Drift struct {
	Added   map[string]string // keys present in current but not in baseline
	Removed map[string]string // keys present in baseline but not in current
	Changed map[string][2]string // key -> [baseline value, current value]
}

// HasDrift returns true if any differences exist between the two snapshots.
func (d Drift) HasDrift() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Take creates a new Snapshot from the given env map and source label.
func Take(source string, env map[string]string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Source:     source,
		Env:        copy,
	}
}

// SaveTo writes a snapshot to the given file path as JSON.
func SaveTo(path string, s Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshotter: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshotter: write %q: %w", path, err)
	}
	return nil
}

// LoadFrom reads a snapshot from the given file path.
func LoadFrom(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshotter: read %q: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshotter: unmarshal: %w", err)
	}
	return s, nil
}

// Compare returns a Drift describing what changed between baseline and current.
func Compare(baseline, current Snapshot) Drift {
	d := Drift{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, cv := range current.Env {
		bv, ok := baseline.Env[k]
		if !ok {
			d.Added[k] = cv
		} else if bv != cv {
			d.Changed[k] = [2]string{bv, cv}
		}
	}
	for k, bv := range baseline.Env {
		if _, ok := current.Env[k]; !ok {
			d.Removed[k] = bv
		}
	}
	return d
}
