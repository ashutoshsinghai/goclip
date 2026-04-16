//go:build darwin

package tray

import "os/exec"

// notifyCopied fires a native macOS notification toast.
// The preview text is passed as a separate argv argument — no shell-escaping needed.
func notifyCopied(preview string) {
	exec.Command("osascript",
		"-e", "on run argv",
		"-e", `display notification (item 1 of argv) with title "goclip" subtitle "Copied to clipboard"`,
		"-e", "end run",
		preview,
	).Start() //nolint:errcheck
}

// openPicker opens a new Terminal window running the interactive TUI picker.
func openPicker() {
	exe := goclipExe()
	exec.Command("osascript",
		"-e", "on run argv",
		"-e", `tell application "Terminal"`,
		"-e", `  activate`,
		"-e", `  do script ((item 1 of argv) & " pick")`,
		"-e", `end tell`,
		"-e", "end run",
		exe,
	).Start() //nolint:errcheck
}
