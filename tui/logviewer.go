package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogViewer struct {
	viewport    viewport.Model
	logs        []string
	maxLogLines int
	writer      *LogViewWriter
}

func NewLogViewer() *LogViewer {
	v := viewport.New(0, 0)
	// v.KeyMap = viewport.KeyMap{
	// 	PageDown: key.NewBinding(
	// 		key.WithKeys("pgdown", "f"),
	// 	),
	// 	PageUp: key.NewBinding(
	// 		key.WithKeys("pgup", "b"),
	// 	),
	// 	HalfPageUp: key.NewBinding(
	// 		key.WithKeys("ctrl+u"),
	// 	),
	// 	HalfPageDown: key.NewBinding(
	// 		key.WithKeys("ctrl+d"),
	// 	),
	// 	Up: key.NewBinding(
	// 		key.WithKeys("up", "k"),
	// 	),
	// 	Down: key.NewBinding(
	// 		key.WithKeys("down", "j"),
	// 	),
	//}
	return &LogViewer{
		viewport:    v,
		logs:        make([]string, 0),
		maxLogLines: 1000, // Keep last 1000 lines
	}
}

func (l *LogViewer) AddLog(log string) {
	// Check if we're at the bottom before adding the new log
	wasAtBottom := l.viewport.AtBottom()

	l.logs = append(l.logs, strings.TrimSpace(log))
	if len(l.logs) > l.maxLogLines {
		l.logs = l.logs[len(l.logs)-l.maxLogLines:]
	}

	l.viewport.SetContent(strings.Join(l.logs, "\n"))

	// Only scroll to bottom if we were already there
	if wasAtBottom {
		l.viewport.GotoBottom()
	}
}

func (l *LogViewer) RefreshContent() {
	if len(l.logs) > 0 {
		l.viewport.SetContent(strings.Join(l.logs, "\n"))
		l.viewport.GotoBottom()
	}
}

func (l *LogViewer) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.viewport.Width = msg.Width
		l.viewport.Height = msg.Height - 2 // Leave room for title
	}
	return cmd
}

func (l *LogViewer) View() string {
	var s strings.Builder

	// Add title and scroll indicator
	scrollPercent := fmt.Sprintf("%3.f%%", l.viewport.ScrollPercent()*100)
	title := fmt.Sprintf("Log Viewer (press ESC to return) - Use ↑/↓, PgUp/PgDn, or mouse wheel to scroll - %s", scrollPercent)
	s.WriteString(dashboardTitleStyle.Render(title))
	s.WriteString("\n")

	// Add scroll bar
	contentHeight := strings.Count(l.viewport.View(), "\n")
	if contentHeight > l.viewport.Height {
		progress := float64(l.viewport.YOffset) / float64(contentHeight-l.viewport.Height)
		barHeight := l.viewport.Height - 2 // Account for title
		position := int(float64(barHeight) * progress)

		// Create vertical scroll bar
		scrollbar := strings.Builder{}
		for i := 0; i < barHeight; i++ {
			if i == position {
				scrollbar.WriteString("█")
			} else {
				scrollbar.WriteString("│")
			}
			scrollbar.WriteString("\n")
		}

		// Join main content and scrollbar
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			l.viewport.View(),
			"  "+scrollbar.String(),
		)
		s.WriteString(content)
	} else {
		s.WriteString(l.viewport.View())
	}

	return s.String()
}

func (l *LogViewer) SetWriter(writer *LogViewWriter) {
	l.writer = writer
}

func (l *LogViewer) RestoreConsoleLogging() {
	if l.writer != nil {
		l.writer.Headless = true
	}
}
