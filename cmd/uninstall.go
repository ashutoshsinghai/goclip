package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Uninstall removes the goclip binary and optionally the history directory.
func Uninstall() {
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Could not find binary path: %v\n", err)
		os.Exit(1)
	}
	binaryPath, _ = filepath.EvalSymlinks(binaryPath)

	historyDir := func() string {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".goclip")
	}()

	fmt.Println("This will remove goclip from your system.")
	fmt.Printf("  Binary:  %s\n", binaryPath)
	fmt.Printf("  History: %s\n", historyDir)
	fmt.Println()

	removeHistory := confirm("Also delete your clipboard history? [y/N]: ", false)

	if !confirm("Proceed with uninstall? [y/N]: ", false) {
		fmt.Println("Aborted.")
		return
	}

	if removeHistory {
		if err := os.RemoveAll(historyDir); err != nil {
			fmt.Printf("Warning: could not remove history: %v\n", err)
		} else {
			fmt.Printf("Removed %s\n", historyDir)
		}
	}

	if err := removeBinary(binaryPath); err != nil {
		fmt.Printf("Could not remove binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed %s\n", binaryPath)
	fmt.Println("goclip uninstalled. Goodbye!")
}

// removeBinary deletes the running binary. On Windows the process holds a lock
// on its own .exe, so we drop a batch script that deletes it after we exit.
func removeBinary(binaryPath string) error {
	if runtime.GOOS != "windows" {
		return os.Remove(binaryPath)
	}

	bat, err := os.CreateTemp("", "goclip-uninstall-*.bat")
	if err != nil {
		return err
	}
	batPath := bat.Name()
	fmt.Fprintf(bat, "@echo off\r\n:loop\r\ndel /f /q %q\r\nif exist %q goto loop\r\ndel /f /q %%~f0\r\n", binaryPath, binaryPath)
	bat.Close()

	cmd := exec.Command("cmd", "/c", "start", "/min", "", batPath)
	return cmd.Start()
}

// confirm prints a prompt and returns true if the user types "y" or "yes".
func confirm(prompt string, defaultYes bool) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}
