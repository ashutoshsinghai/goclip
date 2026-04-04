//go:build !windows

package daemon

import (
	"os"
	"os/exec"
	"syscall"
)

// spawnBackground launches goclip daemon as a detached background process.
// stdout/stderr go to the log file.
func spawnBackground() (int, error) {
	exe, _ := os.Executable()
	logF, err := os.OpenFile(daemonLogFile(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	cmd := exec.Command(exe, "run")
	// Setsid detaches the process from the current terminal session
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout = logF
	cmd.Stderr = logF

	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}

// isProcessAlive sends signal 0 to the PID — no-op but returns an error if dead.
func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// killProcess sends SIGTERM so the daemon can shut down cleanly.
func killProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGTERM)
}
