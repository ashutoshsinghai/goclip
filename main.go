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
	"github.com/ashutoshsinghai/goclip/cmd/tray"
	"github.com/ashutoshsinghai/goclip/internal/style"
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
	case "tray":
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "stop":
				tray.StopTray()
			case "status":
				tray.TrayStatus()
			default:
				fmt.Printf("Unknown tray subcommand %q. Usage: goclip tray [stop|status]\n", os.Args[2])
				os.Exit(1)
			}
		} else {
			tray.StartTray()
		}
	case "tray-run": // internal — spawned by StartTray, not listed in help
		tray.Run()
	case "pick", "ui":
		ui.RunPicker()
	case "list", "ls":
		cmd.ListClips()
	case "search", "find":
		if len(os.Args) < 3 {
			fmt.Println("Usage: goclip search <keyword>")
			os.Exit(1)
		}
		ui.RunPickerWithQuery(strings.Join(os.Args[2:], " "))
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
		force := len(os.Args) > 2 && (os.Args[2] == "--force" || os.Args[2] == "-f")
		cmd.ClearHistory(force)
	case "upgrade":
		cmd.Upgrade(version)
	case "uninstall":
		cmd.Uninstall()
	case "version", "--version", "-v":
		fmt.Println("goclip", version)
	case "help", "--help", "-h":
		printHelp()
	default:
		if suggestion := suggest(os.Args[1]); suggestion != "" {
			fmt.Println(style.Red.Render("Unknown command: ") + style.Dim.Render("\"" + os.Args[1] + "\""))
			fmt.Println(style.Green.Render("✓ Did you mean: ") + style.Bold.Render("goclip " + suggestion) + "?")
		} else {
			fmt.Println(style.Red.Render("Unknown command: ") + style.Dim.Render("\"" + os.Args[1] + "\""))
			fmt.Println(style.Dim.Render("Run ") + style.Bold.Render("goclip help") + style.Dim.Render(" for usage."))
		}
		os.Exit(1)
	}
}

// known is the full list of valid commands used for fuzzy suggestion.
var known = []string{
	"daemon", "stop", "status", "run", "pick", "tray",
	"list", "search", "copy", "pin", "clear",
	"upgrade", "uninstall", "version", "help",
}

// suggest returns the closest known command to input, or "" if nothing is close.
func suggest(input string) string {
	best, bestDist := "", 4 // only suggest if distance ≤ 3
	for _, cmd := range known {
		if d := levenshtein(input, cmd); d < bestDist {
			best, bestDist = cmd, d
		}
	}
	return best
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	row := make([]int, lb+1)
	for j := range row {
		row[j] = j
	}
	for i := 1; i <= la; i++ {
		prev := i
		for j := 1; j <= lb; j++ {
			val := row[j-1] // diagonal
			if a[i-1] != b[j-1] {
				val = 1 + min3(row[j], prev, row[j-1])
			}
			row[j-1] = prev
			prev = val
		}
		row[lb] = prev
	}
	return row[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func printHelp() {
	fmt.Print(`
goclip — Clipboard History Manager

USAGE:
  goclip tray          Start menu bar / system tray app in background
  goclip tray stop     Stop the tray app
  goclip tray status   Show whether the tray app is running
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
  goclip clear --force Skip confirmation prompt (for automation)
  goclip upgrade       Upgrade goclip to the latest version
  goclip uninstall     Remove goclip from your system
  goclip version       Show current version

TYPICAL WORKFLOW:
  1. goclip daemon     # start in background
  2. goclip pick       # open picker whenever you need an old clip
`)
}
