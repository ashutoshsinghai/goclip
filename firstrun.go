package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashutoshsinghai/goclip/internal/autostart"
	"github.com/ashutoshsinghai/goclip/internal/style"
)

// maybeFirstRunSetup runs once per machine, the first time the user invokes
// goclip in a way that implies they actually want to use it. It silently
// enables the login autostart entry so the user never has to remember to
// re-run `goclip daemon` after a reboot, then prints a one-time notice.
//
// First run is detected by absence of ~/.goclip/. Pre-existing users who
// already have a history directory are NOT auto-enrolled — they can opt in
// with `goclip autostart on` if they want it.
func maybeFirstRunSetup() {
	if skipFirstRunSetup() {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := filepath.Join(home, ".goclip")
	if _, err := os.Stat(dir); err == nil {
		return
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	if err := autostart.Enable(); err != nil {
		return
	}
	fmt.Println(style.Green.Render("✓ goclip will start automatically at login.") +
		style.Dim.Render(" Disable any time with 'goclip autostart off'."))
	fmt.Println()
}

// skipFirstRunSetup returns true for commands where running first-run setup
// would surprise the user (version/help) or would be redundant (commands that
// already manage autostart themselves, or internal entry points spawned by
// the parent process which has already run first-run).
func skipFirstRunSetup() bool {
	if len(os.Args) < 2 {
		return false
	}
	switch os.Args[1] {
	case "version", "--version", "-v",
		"help", "--help", "-h",
		"autostart", "uninstall", "install", "upgrade",
		"run", "tray-run":
		return true
	}
	return false
}
