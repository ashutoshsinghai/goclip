package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ashutoshsinghai/goclip/storage"
)

func pidFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".goclip", "daemon.pid")
}

func daemonLogFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".goclip", "daemon.log")
}

// RunDaemon runs the clipboard watcher in the foreground.
// Stdout goes to the terminal (or to the log file when started via `daemon start`).
func RunDaemon() {
	fmt.Println("goclip daemon running — watching your clipboard (Ctrl+C to stop)")
	clips := storage.Load()
	last := ""

	for {
		text, err := clipboard.ReadAll()
		if err == nil && text != "" && text != last {
			last = text
			clips = storage.AddClip(text, clips)
			storage.Save(clips)

			preview := strings.ReplaceAll(text, "\n", "↵")
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			fmt.Printf("[saved] %s\n", preview)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// StartDaemon spawns the daemon as a background process.
func StartDaemon() {
	if pid, alive := readPID(); alive {
		fmt.Printf("Daemon is already running (PID %d)\n", pid)
		return
	}

	pid, err := spawnBackground()
	if err != nil {
		fmt.Printf("Failed to start daemon: %v\n", err)
		os.Exit(1)
	}

	os.WriteFile(pidFile(), []byte(strconv.Itoa(pid)), 0644)
	fmt.Printf("Daemon started (PID %d)\n", pid)
	fmt.Printf("Logs:  %s\n", daemonLogFile())
	fmt.Printf("Stop:  goclip daemon stop\n")
}

// StopDaemon kills the running background daemon.
func StopDaemon() {
	pid, alive := readPID()
	if !alive {
		fmt.Println("Daemon is not running.")
		os.Remove(pidFile()) // clean up stale PID file if any
		return
	}

	if err := killProcess(pid); err != nil {
		fmt.Printf("Failed to stop daemon: %v\n", err)
		os.Exit(1)
	}

	os.Remove(pidFile())
	fmt.Printf("Daemon stopped (PID %d)\n", pid)
}

// DaemonStatus reports whether the background daemon is running.
func DaemonStatus() {
	pid, alive := readPID()
	if alive {
		fmt.Printf("Daemon is running (PID %d)\n", pid)
		fmt.Printf("Logs: %s\n", daemonLogFile())
	} else {
		fmt.Println("Daemon is not running.")
		fmt.Println("Run `goclip daemon start` to start it in the background.")
		fmt.Println("Run `goclip daemon` to start it in the foreground.")
	}
}

// readPID reads the PID file and returns the PID + whether the process is alive.
func readPID() (int, bool) {
	data, err := os.ReadFile(pidFile())
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, false
	}
	return pid, isProcessAlive(pid)
}
