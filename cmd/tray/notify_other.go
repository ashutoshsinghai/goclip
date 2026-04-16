//go:build !darwin && !windows

package tray

func notifyCopied(_ string) {}
func openPicker()            {}
