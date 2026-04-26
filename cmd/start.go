package cmd

import (
	"github.com/ashutoshsinghai/goclip/cmd/daemon"
	"github.com/ashutoshsinghai/goclip/cmd/tray"
)

// Start brings up everything goclip needs in the background:
// the clipboard watcher, plus the menu bar / system tray on platforms that
// support it. Used by `goclip start` and as the payload of the autostart
// login entry. Idempotent — already-running components are left alone.
func Start() {
	daemon.StartDaemon()
	tray.StartTray()
}
