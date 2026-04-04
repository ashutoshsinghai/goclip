package main

// goclip — Clipboard History Manager
//
// This file is intentionally thin — it just routes CLI args to the
// right package. All real logic lives in:
//   storage/   — read/write history to disk
//   ui/        — interactive TUI picker
//   commands/  — daemon, list, copy, clear

import (
	"fmt"
	"os"

	"github.com/ashutoshsinghai/goclip/commands"
	"github.com/ashutoshsinghai/goclip/ui"
)

func main() {
	if len(os.Args) < 2 {
		ui.RunPicker()
		return
	}

	switch os.Args[1] {
	case "daemon":
		commands.RunDaemon()
	case "pick", "ui":
		ui.RunPicker()
	case "list", "ls":
		commands.ListClips()
	case "copy", "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip copy <id>")
			os.Exit(1)
		}
		commands.CopyClip(os.Args[2])
	case "pin":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip pin <id>")
			os.Exit(1)
		}
		commands.PinClip(os.Args[2])
	case "clear":
		commands.ClearHistory()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`
goclip — Clipboard History Manager

USAGE:
  goclip daemon        Start watching your clipboard (keep this running)
  goclip pick          Open interactive TUI picker  ← the fun one
  goclip list          Show history as plain text
  goclip copy <id>     Copy item by ID (non-interactive)
  goclip pin <id>      Pin/unpin an item (pinned items stay at top)
  goclip clear         Wipe all history

TYPICAL WORKFLOW:
  1. goclip daemon     # run once in a background tab
  2. goclip pick       # open picker whenever you need an old clip
`)
}
