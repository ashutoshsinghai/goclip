//go:build !darwin

package tray

import "fmt"

func traySupported() bool { return false }

func Run() {
	fmt.Println("System tray is not supported on Linux.")
	fmt.Println("Use 'goclip daemon' to watch the clipboard in the background.")
}
