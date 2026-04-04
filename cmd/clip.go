package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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
// If force is true, clears without prompting (for automation).
func ClearHistory(force bool) {
	if force {
		storage.Save([]storage.Clip{})
		fmt.Println(style.Yellow.Render("History cleared."))
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(style.Yellow.Render("Are you sure? ") + "This will delete all history. (yes/no): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "yes" || input == "y" {
		storage.Save([]storage.Clip{})
		fmt.Println(style.Yellow.Render("History cleared."))
	} else if input == "no" || input == "n" {
		fmt.Println(style.Dim.Render("Cancelled."))
	} else {
		fmt.Println(style.Red.Render("Invalid input. ") + "Please enter 'yes' or 'no'.")
		ClearHistory(false)
	}
}
