package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

// Uninstall removes the goclip binary and optionally the history directory.
func Uninstall() {
	// Find where this binary lives
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

	// Ask about history
	removeHistory := confirm("Also delete your clipboard history? [y/N]: ", false)

	// Final confirmation
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

	// Remove the binary last — once this runs, we're gone
	if err := os.Remove(binaryPath); err != nil {
		fmt.Printf("Could not remove binary (try with sudo): %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed %s\n", binaryPath)
	fmt.Println("goclip uninstalled. Goodbye!")
}

// confirm prints a prompt and returns true if the user types "y" or "yes".
// defaultYes controls what Enter alone means.
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
