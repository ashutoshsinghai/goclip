//go:build windows

package autostart

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey   = `Software\Microsoft\Windows\CurrentVersion\Run`
	runValue = "goclip"
)

func enable() error {
	exe, err := goclipBinary()
	if err != nil {
		return err
	}

	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	// Quote the path so spaces in the install location don't break parsing.
	cmd := `"` + exe + `" start`
	return k.SetStringValue(runValue, cmd)
}

func disable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}
	defer k.Close()

	if err := k.DeleteValue(runValue); err != nil {
		// "value does not exist" → already disabled, treat as success.
		if strings.Contains(err.Error(), "cannot find") {
			return nil
		}
		return err
	}
	return nil
}

func isEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	_, _, err = k.GetStringValue(runValue)
	return err == nil
}
