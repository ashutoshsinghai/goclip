// Package autostart manages goclip's "start at login" entry on the host OS.
//
// Each platform has its own native mechanism:
//
//	macOS    ~/Library/LaunchAgents/com.goclip.plist     (loaded via launchctl)
//	Linux    ~/.config/autostart/goclip.desktop          (XDG autostart)
//	Windows  HKCU\Software\Microsoft\Windows\CurrentVersion\Run\goclip
//
// Enable/Disable are idempotent — calling Enable when already enabled is a
// no-op, and Disable on a system without an entry is a no-op too. They never
// fail loudly: autostart is a quality-of-life feature, not a hard requirement,
// and the binary still works without it.
package autostart

// Enable installs the login entry so goclip starts automatically on login.
// Returns nil if already enabled.
func Enable() error { return enable() }

// Disable removes the login entry. Returns nil if it wasn't installed.
func Disable() error { return disable() }

// IsEnabled reports whether the login entry currently exists.
func IsEnabled() bool { return isEnabled() }
