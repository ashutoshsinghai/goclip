package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/ashutoshsinghai/goclip/storage"
)

// CopyClip puts a historical clip back on the clipboard by its ID.
func CopyClip(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(red.Render("Error: ") + "ID must be a number, e.g. goclip copy 3")
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
			fmt.Println(green.Render("Copied: ") + preview)
			return
		}
	}
	fmt.Println(red.Render(fmt.Sprintf("No clip found with ID %d", id)))
}

// PinClip toggles the pin on a clip by ID.
func PinClip(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(red.Render("Error: ") + "ID must be a number, e.g. goclip pin 3")
		os.Exit(1)
	}
	clips := storage.Load()
	clips, pinned := storage.TogglePin(id, clips)
	if pinned {
		fmt.Println(yellow.Render("★ Pinned ") + dim.Render(fmt.Sprintf("clip #%d", id)))
	} else {
		fmt.Println(dim.Render(fmt.Sprintf("Unpinned clip #%d", id)))
	}
	storage.Save(clips)
}

// ClearHistory wipes all saved clipboard history.
func ClearHistory() {
	storage.Save([]storage.Clip{})
	fmt.Println(yellow.Render("History cleared."))
}
