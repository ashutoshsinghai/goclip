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
	dir := filepath.Join(home, ".goclip")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "daemon.log")
}

// RunDaemon runs the clipboard watcher in the foreground.
// Stdout goes to the terminal (or to the log file when started via `daemon start`).
func RunDaemon() {
	fmt.Println("goclip daemon running — watching your clipboard (Ctrl+C to stop)")

	// Startup check: verify clipboard is accessible before entering the loop.
	// On Linux this fails if xclip/xsel/wl-clipboard is missing or DISPLAY is unset.
	if _, err := clipboard.ReadAll(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: clipboard not accessible:", err)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "On Linux, install a clipboard utility and ensure DISPLAY is set:")
		fmt.Fprintln(os.Stderr, "  sudo apt install xclip        # X11")
		fmt.Fprintln(os.Stderr, "  sudo apt install wl-clipboard # Wayland")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "If connecting via SSH, use: ssh -X user@host")
		os.Exit(1)
	}

	clips := storage.Load()
	last := ""
	var lastErrTime time.Time

	for {
		text, err := clipboard.ReadAll()
		if err != nil {
			// Log at most once per minute to avoid flooding the log file.
			if time.Since(lastErrTime) > time.Minute {
				fmt.Fprintln(os.Stderr, "[error] clipboard read failed:", err)
				lastErrTime = time.Now()
			}
		} else if text != "" && text != last {
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
