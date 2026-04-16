//go:build !windows

package tray

import (
	"os"
	"os/exec"
	"syscall"
)

// spawnTray launches "goclip tray-run" as a detached background process.
// stdout/stderr go to the tray log file.
func spawnTray() (int, error) {
	exe, _ := os.Executable()
	logF, err := os.OpenFile(trayLogFile(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	cmd := exec.Command(exe, "tray-run")
	// Setpgid puts the child in a new process group so it won't receive
	// Ctrl-C from the parent terminal, but it stays in the same *session*
	// so macOS grants it a window-server connection (unlike Setsid which
	// creates a new session and loses the window-server context entirely).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = logF
	cmd.Stderr = logF

	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}

func isTrayProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

func killTrayProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGTERM)
}
