//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

func desktopPath() string {
	cfg := os.Getenv("XDG_CONFIG_HOME")
	if cfg == "" {
		home, _ := os.UserHomeDir()
		cfg = filepath.Join(home, ".config")
	}
	return filepath.Join(cfg, "autostart", "goclip.desktop")
}

func enable() error {
	exe, err := goclipBinary()
	if err != nil {
		return err
	}

	path := desktopPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=goclip
Comment=Clipboard history daemon
Exec=%s start
Terminal=false
X-GNOME-Autostart-enabled=true
`, exe)

	return os.WriteFile(path, []byte(content), 0644)
}

func disable() error {
	if err := os.Remove(desktopPath()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func isEnabled() bool {
	_, err := os.Stat(desktopPath())
	return err == nil
}
