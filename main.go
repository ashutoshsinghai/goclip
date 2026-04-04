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
	"strings"

	"github.com/ashutoshsinghai/goclip/commands"
	"github.com/ashutoshsinghai/goclip/ui"
)

// version is set at build time via goreleaser ldflags.
// Falls back to "dev" when built locally with `go build`.
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		ui.RunPicker()
		return
	}

	switch os.Args[1] {
	case "daemon":
		if len(os.Args) < 3 {
			commands.RunDaemon() // foreground
			return
		}
		switch os.Args[2] {
		case "start":
			commands.StartDaemon()
		case "stop":
			commands.StopDaemon()
		case "status":
			commands.DaemonStatus()
		default:
			fmt.Printf("Unknown daemon subcommand: %s\n", os.Args[2])
			fmt.Println("Usage: goclip daemon [start|stop|status]")
			os.Exit(1)
		}
	case "pick", "ui":
		ui.RunPicker()
	case "list", "ls":
		commands.ListClips()
	case "search", "find":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip search <keyword>")
			os.Exit(1)
		}
		commands.SearchClips(strings.Join(os.Args[2:], " "))
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
	case "upgrade":
		commands.Upgrade(version)
	case "uninstall":
		commands.Uninstall()
	case "version", "--version", "-v":
		fmt.Println("goclip", version)
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
  goclip daemon        Start watching clipboard in foreground
  goclip daemon start  Start daemon in background
  goclip daemon stop   Stop background daemon
  goclip daemon status Show whether daemon is running
  goclip pick          Open interactive TUI picker  ← the fun one
  goclip list          Show history as plain text
  goclip search <kw>   Search history by keyword
  goclip copy <id>     Copy item by ID (non-interactive)
  goclip pin <id>      Pin/unpin an item (pinned items stay at top)
  goclip clear         Wipe all history
  goclip upgrade       Upgrade goclip to the latest version
  goclip uninstall     Remove goclip from your system
  goclip version       Show current version

TYPICAL WORKFLOW:
  1. goclip daemon     # run once in a background tab
  2. goclip pick       # open picker whenever you need an old clip
`)
}
