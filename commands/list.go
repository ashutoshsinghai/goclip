package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ashutoshsinghai/goclip/storage"
)

// ListClips prints all clipboard history as a plain text table.
func ListClips() {
	clips := storage.Load()
	if len(clips) == 0 {
		fmt.Println(yellow.Render("No clipboard history yet.") + dim.Render(" Run: goclip daemon"))
		return
	}
	fmt.Println()
	fmt.Printf("%-5s  %-17s  %s\n",
		bold.Render("ID"),
		bold.Render("TIME"),
		bold.Render("CONTENT"),
	)
	fmt.Println(dim.Render(strings.Repeat("─", 80)))
	for _, c := range clips {
		preview := strings.ReplaceAll(c.Content, "\n", "↵")
		if len(preview) > 55 {
			preview = preview[:55] + "..."
		}
		pin := "  "
		if c.Pinned {
			pin = yellow.Render("★ ")
		}
		fmt.Printf("%-5d  %-17s  %s%s\n",
			c.ID,
			dim.Render(c.CopiedAt.Format("Jan 02 15:04:05")),
			pin,
			preview,
		)
	}
	fmt.Println()
}

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

// ClearHistory wipes all saved clipboard history.
func ClearHistory() {
	storage.Save([]storage.Clip{})
	fmt.Println(yellow.Render("History cleared."))
}

// SearchClips prints all clips whose content contains the query.
func SearchClips(query string) {
	clips := storage.Load()
	q := strings.ToLower(query)

	var matches []storage.Clip
	for _, c := range clips {
		if strings.Contains(strings.ToLower(c.Content), q) {
			matches = append(matches, c)
		}
	}

	if len(matches) == 0 {
		fmt.Println(yellow.Render("No results for ") + bold.Render(fmt.Sprintf("%q", query)))
		return
	}

	fmt.Println()
	fmt.Printf("%-5s  %-17s  %s\n",
		bold.Render("ID"),
		bold.Render("TIME"),
		bold.Render("CONTENT"),
	)
	fmt.Println(dim.Render(strings.Repeat("─", 80)))
	for _, c := range matches {
		preview := strings.ReplaceAll(c.Content, "\n", "↵")
		if len(preview) > 55 {
			preview = preview[:55] + "..."
		}
		pin := "  "
		if c.Pinned {
			pin = yellow.Render("★ ")
		}
		fmt.Printf("%-5d  %-17s  %s%s\n",
			c.ID,
			dim.Render(c.CopiedAt.Format("Jan 02 15:04:05")),
			pin,
			preview,
		)
	}
	fmt.Printf("\n%s\n\n", dim.Render(fmt.Sprintf("%d result(s) for %q", len(matches), query)))
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
