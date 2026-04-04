package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/ashutoshsinghai/goclip/internal/storage"
	"github.com/ashutoshsinghai/goclip/internal/style"
)

// CopyClip puts a historical clip back on the clipboard by its ID.
func CopyClip(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(style.Red.Render("Error: ") + "ID must be a number, e.g. goclip copy 3")
		os.Exit(1)
	}
	clips := storage.Load()
	for _, c := range clips {
		if c.ID == id {
			clipboard.WriteAll(c.Content)
			preview := c.Content
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			fmt.Println(style.Green.Render("Copied: ") + preview)
			return
		}
	}
	fmt.Println(style.Red.Render(fmt.Sprintf("No clip found with ID %d", id)))
}

// PinClip toggles the pin on a clip by ID.
func PinClip(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(style.Red.Render("Error: ") + "ID must be a number, e.g. goclip pin 3")
		os.Exit(1)
	}
	clips := storage.Load()
	clips, pinned := storage.TogglePin(id, clips)
	if pinned {
		fmt.Println(style.Yellow.Render("★ Pinned ") + style.Dim.Render(fmt.Sprintf("clip #%d", id)))
	} else {
		fmt.Println(style.Dim.Render(fmt.Sprintf("Unpinned clip #%d", id)))
	}
	storage.Save(clips)
}

// ClearHistory wipes all saved clipboard history.
func ClearHistory() {
	storage.Save([]storage.Clip{})
	fmt.Println(style.Yellow.Render("History cleared."))
}
