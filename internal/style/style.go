package style

import "github.com/charmbracelet/lipgloss"

// Shared terminal styles used across cmd/ and cmd/daemon/.
var (
	Green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF88")).Bold(true)
	Red    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	Dim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	Bold   = lipgloss.NewStyle().Bold(true)
)
