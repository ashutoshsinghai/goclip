// Package commands contains the implementations for each CLI subcommand.
package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ashutoshsinghai/goclip/storage"
)

// RunDaemon polls the clipboard every 500ms and saves new entries to disk.
// It runs forever — the user stops it with Ctrl+C.
func RunDaemon() {
	fmt.Println("goclip daemon running — watching your clipboard (Ctrl+C to stop)")

	clips := storage.Load()
	last := ""

	for {
		text, err := clipboard.ReadAll()
		if err == nil && text != "" && text != last {
			last = text
			clips = storage.AddClip(text, clips)
			storage.Save(clips)

			preview := strings.ReplaceAll(text, "\n", "↵")
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			fmt.Printf("[saved] %s\n", preview)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
