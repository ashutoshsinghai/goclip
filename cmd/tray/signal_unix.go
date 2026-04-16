//go:build !windows

package tray

import (
	"os/signal"
	"syscall"
)

// ignoreSighup prevents the tray process from being killed when the terminal
// that launched it closes. Without this, closing the terminal sends SIGHUP to
// all processes in the session, which would take down the menu bar app.
func ignoreSighup() {
	signal.Ignore(syscall.SIGHUP)
}
