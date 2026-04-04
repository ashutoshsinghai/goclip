package main

// goclip — Clipboard History Manager
//
// This file is intentionally thin — it just routes CLI args to the
// right package. All real logic lives in:
//   internal/storage/  — read/write history to disk
//   ui/                — interactive TUI picker
//   cmd/               — history, clip, upgrade, uninstall
//   cmd/daemon/        — daemon start/stop/status/run

import (
	"fmt"
	"os"
	"strings"

	"github.com/ashutoshsinghai/goclip/cmd"
	"github.com/ashutoshsinghai/goclip/cmd/daemon"
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
	case "run":
		daemon.RunDaemon()
	case "daemon":
		daemon.StartDaemon()
	case "stop":
		daemon.StopDaemon()
	case "status":
		daemon.DaemonStatus()
	case "pick", "ui":
		ui.RunPicker()
	case "list", "ls":
		cmd.ListClips()
	case "search", "find":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip search <keyword>")
			os.Exit(1)
		}
		cmd.SearchClips(strings.Join(os.Args[2:], " "))
	case "copy", "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip copy <id>")
			os.Exit(1)
		}
		cmd.CopyClip(os.Args[2])
	case "pin":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip pin <id>")
			os.Exit(1)
		}
		cmd.PinClip(os.Args[2])
	case "clear":
		cmd.ClearHistory()
	case "upgrade":
		cmd.Upgrade(version)
	case "uninstall":
		cmd.Uninstall()
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
  goclip daemon        Start clipboard watcher in background
  goclip stop          Stop background daemon
  goclip status        Show whether daemon is running
  goclip run           Run clipboard watcher in foreground
  goclip pick          Open interactive TUI picker
  goclip list          Show history as plain text
  goclip search <kw>   Search history by keyword
  goclip copy <id>     Copy item by ID (non-interactive)
  goclip pin <id>      Pin/unpin an item (pinned items stay at top)
  goclip clear         Wipe all history
  goclip upgrade       Upgrade goclip to the latest version
  goclip uninstall     Remove goclip from your system
  goclip version       Show current version

TYPICAL WORKFLOW:
  1. goclip daemon     # start in background
  2. goclip pick       # open picker whenever you need an old clip
`)
}
