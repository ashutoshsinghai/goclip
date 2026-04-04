package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ashutoshsinghai/goclip/internal/storage"
	"github.com/ashutoshsinghai/goclip/internal/style"
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
		fmt.Println(style.Yellow.Render("Daemon is already running") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
		return
	}

	pid, err := spawnBackground()
	if err != nil {
		fmt.Println(style.Red.Render("Failed to start daemon: ") + err.Error())
		os.Exit(1)
	}

	os.WriteFile(pidFile(), []byte(strconv.Itoa(pid)), 0644)
	fmt.Println(style.Green.Render("Daemon started") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
	fmt.Println(style.Dim.Render("Logs:  " + daemonLogFile()))
	fmt.Println(style.Dim.Render("Stop:  goclip stop"))
}

// StopDaemon kills the running background daemon.
func StopDaemon() {
	pid, alive := readPID()
	if !alive {
		fmt.Println(style.Yellow.Render("Daemon is not running."))
		os.Remove(pidFile())
		return
	}

	if err := killProcess(pid); err != nil {
		fmt.Println(style.Red.Render("Failed to stop daemon: ") + err.Error())
		os.Exit(1)
	}

	os.Remove(pidFile())
	fmt.Println(style.Yellow.Render("Daemon stopped") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
}

// DaemonStatus reports whether the background daemon is running.
func DaemonStatus() {
	pid, alive := readPID()
	if alive {
		fmt.Println(style.Green.Render("● Daemon is running") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
		fmt.Println(style.Dim.Render("  Logs: " + daemonLogFile()))
		fmt.Println(style.Dim.Render("  Stop: goclip stop"))
	} else {
		fmt.Println(style.Red.Render("○ Daemon is not running"))
		fmt.Println(style.Dim.Render("  Start: goclip daemon"))
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
