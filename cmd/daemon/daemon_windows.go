//go:build windows

package daemon

import (
	"os"
	"os/exec"
	"syscall"
)

const detachedProcess = 0x00000008

// spawnBackground launches goclip daemon as a hidden background process on Windows.
func spawnBackground() (int, error) {
	exe, _ := os.Executable()
	logF, err := os.OpenFile(daemonLogFile(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	cmd := exec.Command(exe, "run")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | detachedProcess,
		HideWindow:    true,
	}
	cmd.Stdout = logF
	cmd.Stderr = logF

	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}

// isProcessAlive checks if a process is running on Windows via OpenProcess.
func isProcessAlive(pid int) bool {
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	syscall.CloseHandle(handle)
	return true
}

// killProcess forcefully kills the process on Windows.
func killProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
