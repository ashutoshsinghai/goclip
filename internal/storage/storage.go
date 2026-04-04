// Package storage handles reading and writing clipboard history to disk.
package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const MaxHistory = 100

// Clip represents a single clipboard entry.
type Clip struct {
	ID       int       `json:"id"`
	Content  string    `json:"content"`
	CopiedAt time.Time `json:"copied_at"`
	Pinned   bool      `json:"pinned,omitempty"`
}

// HistoryFile returns the path to ~/.goclip/history.json,
// creating the directory if it doesn't exist yet.
func HistoryFile() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".goclip")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "history.json")
}

// Load reads clipboard history from disk.
// Returns an empty slice if the file doesn't exist yet.
func Load() []Clip {
	data, err := os.ReadFile(HistoryFile())
	if err != nil {
		return []Clip{}
	}
	var clips []Clip
	if err := json.Unmarshal(data, &clips); err != nil {
		// Return empty slice if JSON is corrupted, don't panic
		return []Clip{}
	}
	// Ensure we never return nil, always an empty slice
	if clips == nil {
		return []Clip{}
	}
	return clips
}

// Save writes the clipboard history back to disk as JSON.
func Save(clips []Clip) {
	data, _ := json.MarshalIndent(clips, "", "  ")
	os.WriteFile(HistoryFile(), data, 0644)
}

// AddClip prepends a new clip to the list.
// Skips if identical to the most recent entry. Caps at MaxHistory.
func AddClip(content string, clips []Clip) []Clip {
	if len(clips) > 0 && clips[0].Content == content {
		return clips
	}

	newID := 1
	if len(clips) > 0 {
		newID = clips[0].ID + 1
	}

	clip := Clip{ID: newID, Content: content, CopiedAt: time.Now()}
	clips = append([]Clip{clip}, clips...)

	// Trim unpinned items to MaxHistory. Pinned items are never removed.
	unpinned := 0
	var trimmed []Clip
	for _, c := range clips {
		if c.Pinned {
			trimmed = append(trimmed, c)
		} else if unpinned < MaxHistory {
			trimmed = append(trimmed, c)
			unpinned++
		}
	}
	return trimmed
}

// TogglePin flips the Pinned flag on the clip with the given ID.
func TogglePin(id int, clips []Clip) ([]Clip, bool) {
	for i, c := range clips {
		if c.ID == id {
			clips[i].Pinned = !c.Pinned
			return clips, clips[i].Pinned
		}
	}
	return clips, false
}

// Sorted returns clips with pinned items first, preserving relative order within each group.
func Sorted(clips []Clip) []Clip {
	var pinned, rest []Clip
	for _, c := range clips {
		if c.Pinned {
			pinned = append(pinned, c)
		} else {
			rest = append(rest, c)
		}
	}
	return append(pinned, rest...)
}
