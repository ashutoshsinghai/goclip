// Package ui contains the bubbletea TUI picker.
package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ashutoshsinghai/goclip/internal/storage"
)

// -----------------------------------------------------------------------
// Styles  (lipgloss = CSS for terminals)
// -----------------------------------------------------------------------

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7B61FF")).
			Padding(0, 2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7B61FF"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	previewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B61FF")).
			Padding(1, 2).
			Foreground(lipgloss.Color("#DDDDDD")).
			Width(70).
			Height(6)

	searchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF88")).
			Bold(true)

	pinnedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)
)

// -----------------------------------------------------------------------
// Model
// -----------------------------------------------------------------------

// model holds all the state for the TUI.
type model struct {
	clips      []storage.Clip // full history
	filtered   []storage.Clip // after search filter
	cursor     int
	search     string
	searching  bool
	copied     bool
	copiedText string
	quitting   bool
	width      int // terminal width
	height     int // terminal height
}

func newModel() model {
	clips := storage.Sorted(storage.Load())
	return model{
		clips:    clips,
		filtered: clips,
		width:    80,  // default terminal width
		height:   24,  // default terminal height
	}
}

// Init runs once on startup — nothing to do here.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles every keypress and returns the new model state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.TickMsg:
		// Auto-quit 500ms after copying
		if m.copied {
			return m, tea.Quit
		}
	case tea.KeyMsg:
		// After copying, any key quits immediately
		if m.copied {
			return m, tea.Quit
		}
		if m.searching {
			return m.handleSearchKey(msg)
		}
		return m.handleBrowseKey(msg)
	}
	return m, nil
}

func (m model) handleBrowseKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case "enter":
		if len(m.filtered) > 0 {
			clipboard.WriteAll(m.filtered[m.cursor].Content)
			m.copied = true
			m.copiedText = m.filtered[m.cursor].Content
			// Auto-quit after 500ms
			return m, tea.Tick(time.Duration(500)*time.Millisecond, func(t time.Time) tea.Msg {
				return tea.TickMsg{}
			})
		}
	case "/":
		m.searching = true
		m.search = ""
		m.filtered = m.clips
		m.cursor = 0

	case "p":
		// Toggle pin on the selected clip, save immediately
		if len(m.filtered) > 0 {
			selected := m.filtered[m.cursor]
			allClips := storage.Load()
			allClips, _ = storage.TogglePin(selected.ID, allClips)
			storage.Save(allClips)
			// Reload sorted list, keep cursor in place
			sorted := storage.Sorted(allClips)
			m.clips = sorted
			m.filtered = filterClips(sorted, m.search)
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
		}
	}
	return m, nil
}

func (m model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searching = false
		m.search = ""
		m.filtered = m.clips
		m.cursor = 0
	case "enter":
		m.searching = false
		if len(m.filtered) > 0 {
			clipboard.WriteAll(m.filtered[m.cursor].Content)
			m.copied = true
			m.copiedText = m.filtered[m.cursor].Content
		}
	case "backspace":
		if len(m.search) > 0 {
			m.search = m.search[:len(m.search)-1]
			m.filtered = filterClips(m.clips, m.search)
			m.cursor = 0
		}
	default:
		if len(msg.String()) == 1 {
			m.search += msg.String()
			m.filtered = filterClips(m.clips, m.search)
			m.cursor = 0
		}
	}
	return m, nil
}

func filterClips(clips []storage.Clip, query string) []storage.Clip {
	if query == "" {
		return clips
	}
	q := strings.ToLower(query)
	var result []storage.Clip
	for _, c := range clips {
		if strings.Contains(strings.ToLower(c.Content), q) {
			result = append(result, c)
		}
	}
	return result
}

// View renders the whole UI as a string on every update.
func (m model) View() string {
	if m.quitting && !m.copied {
		return ""
	}

	if m.copied {
		preview := m.copiedText
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		return "\n  " + successStyle.Render("Copied: "+preview) + "\n\n"
	}

	var b strings.Builder

	b.WriteString("  " + titleStyle.Render(" goclip — Clipboard History ") + "\n")

	if len(m.clips) == 0 {
		b.WriteString("  " + dimStyle.Render("No history yet. Run: goclip daemon") + "\n")
		b.WriteString("  " + helpStyle.Render("q to quit") + "\n")
		return b.String()
	}

	if len(m.filtered) == 0 {
		b.WriteString("  " + dimStyle.Render("No matches found.") + "\n")
		b.WriteString("  " + helpStyle.Render("esc clear search") + "\n")
		return b.String()
	}

	// Ensure cursor is in bounds
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}

	// Show 12 items per page (use more terminal space)
	pageSize := 12
	pageStart := (m.cursor / pageSize) * pageSize

	b.WriteString("\n")
	// Render visible items
	for i := pageStart; i < pageStart+pageSize && i < len(m.filtered); i++ {
		c := m.filtered[i]
		preview := strings.ReplaceAll(c.Content, "\n", " ")
		if len(preview) > 45 {
			preview = preview[:45] + "..."
		}

		pin := "  "
		if c.Pinned {
			pin = pinnedStyle.Render("★ ")
		}

		timeStr := dimStyle.Render(c.CopiedAt.Format("Jan 02 15:04"))
		line := fmt.Sprintf("%s%s  %s", pin, timeStr, preview)

		if i == m.cursor {
			b.WriteString("  " + selectedStyle.Render("▶ "+line) + "\n")
		} else {
			b.WriteString("  " + dimStyle.Render(line) + "\n")
		}
	}

	// Simple preview - just show first line
	if m.cursor >= 0 && m.cursor < len(m.filtered) {
		preview := m.filtered[m.cursor].Content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		preview = strings.ReplaceAll(preview, "\n", " ↵ ")
		b.WriteString("\n  " + dimStyle.Render("Preview: "+preview) + "\n")
	}

	// Help
	b.WriteString("\n  " + helpStyle.Render("↑/↓ navigate   enter copy   p pin   / search   q quit") + "\n")

	return b.String()
}

// -----------------------------------------------------------------------
// Public entry point
// -----------------------------------------------------------------------

// RunPicker opens the interactive TUI clipboard picker.
func RunPicker() {
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithOutput(os.Stderr), // Use stderr for better compatibility
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Picker error: %v\n", err)
	}
}
