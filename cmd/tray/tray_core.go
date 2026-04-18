package tray

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ashutoshsinghai/goclip/internal/style"
)

func trayDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".goclip")
	os.MkdirAll(dir, 0755)
	return dir
}

func trayPidFile() string { return filepath.Join(trayDir(), "tray.pid") }
func trayLogFile() string { return filepath.Join(trayDir(), "tray.log") }

func goclipExe() string {
	if path, err := exec.LookPath("goclip"); err == nil {
		return path
	}
	exe, _ := os.Executable()
	return exe
}

func StartTray() {
	if !traySupported() {
		fmt.Println(style.Yellow.Render("System tray is not supported on this platform."))
		fmt.Println(style.Dim.Render("Use 'goclip daemon' to watch the clipboard in the background."))
		return
	}
	if pid, alive := readTrayPID(); alive {
		fmt.Println(style.Yellow.Render("Tray is already running") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
		return
	}
	pid, err := spawnTray()
	if err != nil {
		fmt.Println(style.Red.Render("Failed to start tray: ") + err.Error())
		os.Exit(1)
	}
	os.WriteFile(trayPidFile(), []byte(strconv.Itoa(pid)), 0644)
	fmt.Println(style.Green.Render("Tray started") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
	fmt.Println(style.Dim.Render("Stop:   goclip tray stop"))
	fmt.Println(style.Dim.Render("Status: goclip tray status"))
}

func StopTray() {
	pid, alive := readTrayPID()
	if !alive {
		fmt.Println(style.Yellow.Render("Tray is not running."))
		os.Remove(trayPidFile())
		return
	}
	if err := killTrayProcess(pid); err != nil {
		fmt.Println(style.Red.Render("Failed to stop tray: ") + err.Error())
		os.Exit(1)
	}
	os.Remove(trayPidFile())
	fmt.Println(style.Yellow.Render("Tray stopped") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
}

func TrayStatus() {
	pid, alive := readTrayPID()
	if alive {
		fmt.Println(style.Green.Render("● Tray is running") + style.Dim.Render(fmt.Sprintf(" (PID %d)", pid)))
		fmt.Println(style.Dim.Render("  Stop: goclip tray stop"))
	} else {
		fmt.Println(style.Red.Render("○ Tray is not running"))
		fmt.Println(style.Dim.Render("  Start: goclip tray"))
	}
}

func readTrayPID() (int, bool) {
	data, err := os.ReadFile(trayPidFile())
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, false
	}
	return pid, isTrayProcessAlive(pid)
}
