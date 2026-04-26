package cmd

import (
	"fmt"
	"os"

	"github.com/ashutoshsinghai/goclip/internal/autostart"
	"github.com/ashutoshsinghai/goclip/internal/style"
)

// Autostart manages the "start at login" entry. Subcommands: on / off / status.
// An empty subcommand is treated as "status" so that bare `goclip autostart`
// just prints whether it's currently enabled.
func Autostart(sub string) {
	switch sub {
	case "", "status":
		if autostart.IsEnabled() {
			fmt.Println(style.Green.Render("● Autostart is enabled") + style.Dim.Render(" — goclip will start at login"))
		} else {
			fmt.Println(style.Red.Render("○ Autostart is disabled") + style.Dim.Render(" — run 'goclip autostart on' to enable"))
		}
	case "on", "enable":
		if err := autostart.Enable(); err != nil {
			fmt.Println(style.Red.Render("Failed to enable autostart: ") + err.Error())
			os.Exit(1)
		}
		fmt.Println(style.Green.Render("✓ Autostart enabled") + style.Dim.Render(" — goclip will start at login"))
	case "off", "disable":
		if err := autostart.Disable(); err != nil {
			fmt.Println(style.Red.Render("Failed to disable autostart: ") + err.Error())
			os.Exit(1)
		}
		fmt.Println(style.Yellow.Render("✓ Autostart disabled"))
	default:
		fmt.Printf("Unknown autostart subcommand %q. Usage: goclip autostart [on|off|status]\n", sub)
		os.Exit(1)
	}
}
