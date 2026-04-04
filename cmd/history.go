package cmd

import (
	"fmt"
	"strings"

	"github.com/ashutoshsinghai/goclip/internal/storage"
	"github.com/ashutoshsinghai/goclip/internal/style"
)

// ListClips prints all clipboard history as a plain text table.
func ListClips() {
	clips := storage.Load()
	if len(clips) == 0 {
		fmt.Println(style.Yellow.Render("No clipboard history yet.") + style.Dim.Render(" Run: goclip daemon"))
		return
	}
	fmt.Println()
	fmt.Printf("%-5s  %-17s  %s\n",
		style.Bold.Render("ID"),
		style.Bold.Render("TIME"),
		style.Bold.Render("CONTENT"),
	)
	fmt.Println(style.Dim.Render(strings.Repeat("─", 80)))
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
		fmt.Println(style.Yellow.Render("No results for ") + style.Bold.Render(fmt.Sprintf("%q", query)))
		return
	}

	fmt.Println()
	fmt.Printf("%-5s  %-17s  %s\n",
		style.Bold.Render("ID"),
		style.Bold.Render("TIME"),
		style.Bold.Render("CONTENT"),
	)
	fmt.Println(style.Dim.Render(strings.Repeat("─", 80)))
	for _, c := range matches {
		printClipRow(c)
	}
	fmt.Printf("\n%s\n\n", style.Dim.Render(fmt.Sprintf("%d result(s) for %q", len(matches), query)))
}

// printClipRow prints a single clip as a table row.
func printClipRow(c storage.Clip) {
	preview := strings.ReplaceAll(c.Content, "\n", "↵")
	if len(preview) > 55 {
		preview = preview[:55] + "..."
	}
	pin := "  "
	if c.Pinned {
		pin = style.Yellow.Render("★ ")
	}
	fmt.Printf("%-5d  %-17s  %s%s\n",
		c.ID,
		style.Dim.Render(c.CopiedAt.Format("Jan 02 15:04:05")),
		pin,
		preview,
	)
}
