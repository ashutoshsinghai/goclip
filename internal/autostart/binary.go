package autostart

import (
	"os"
	"os/exec"
	"path/filepath"
)

// goclipBinary returns the absolute path to the goclip binary that the
// autostart entry should invoke. Prefers the binary on $PATH (so updates via
// "goclip install" don't leave the autostart entry pointing at a stale copy
// when the user has it installed in a stable location), but falls back to the
// running executable.
func goclipBinary() (string, error) {
	if path, err := exec.LookPath("goclip"); err == nil {
		abs, err := filepath.Abs(path)
		if err == nil {
			return filepath.EvalSymlinks(abs)
		}
	}
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		return resolved, nil
	}
	return exe, nil
}
