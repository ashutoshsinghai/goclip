//go:build windows

package tray

import (
	"os/exec"
	"time"

	"fyne.io/systray"
)

func openPicker() {
	exe := goclipExe()
	exec.Command("cmd", "/C", "start", "cmd", "/K", exe+" pick").Start() //nolint:errcheck
}

func notifyCopied(preview string) {
	msg := "✓ " + preview
	runes := []rune(msg)
	if len(runes) > 40 {
		msg = string(runes[:39]) + "…"
	}
	systray.SetTooltip(msg)
	time.AfterFunc(2*time.Second, func() {
		systray.SetTooltip("goclip — Clipboard History")
	})
}
