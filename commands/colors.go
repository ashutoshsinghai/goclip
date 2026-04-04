package commands

import "github.com/charmbracelet/lipgloss"

// Shared styles for CLI output across all commands.
var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF88")).Bold(true)
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	dim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	bold   = lipgloss.NewStyle().Bold(true)
)
