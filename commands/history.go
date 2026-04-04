package commands

import (
	"fmt"
	"strings"

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
		printClipRow(c)
	}
	fmt.Println()
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
		printClipRow(c)
	}
	fmt.Printf("\n%s\n\n", dim.Render(fmt.Sprintf("%d result(s) for %q", len(matches), query)))
}

// printClipRow prints a single clip as a table row.
func printClipRow(c storage.Clip) {
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
