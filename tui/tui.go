package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

type IValidatorManager interface {
	GetCurrentReward() float64
	GetTotalTasks() int64
	GetClaims() int64
	GetSolutions() int64
	GetCommitments() int64
	GetValidatorInfo() string
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	focusedModelStyle = lipgloss.NewStyle().
				Width(20).
				Height(20).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type UIinstance struct {
	content   string
	ready     bool
	viewport  viewport.Model
	validator IValidatorManager
}

func (m UIinstance) Init() tea.Cmd {

	return nil
}

func (m UIinstance) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case string:
		//log.Printf("msg: %s", msg)
		// Handle log messages.
		m.content += msg
		m.viewport.SetYOffset(m.viewport.YOffset + 1)
		m.viewport.SetContent(m.content)
		//m.viewport.GotoBottom()

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width/2, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width / 2
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m UIinstance) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var s string

	reward := m.validator.GetCurrentReward()
	info := m.validator.GetValidatorInfo()

	infoBox := focusedModelStyle.MaxWidth(20).Render(fmt.Sprintf("reward: %f\n%s\n", reward, info))

	s += m.headerView() + "\n"
	s += infoBox //lipgloss.JoinHorizontal(lipgloss.Top, m.viewport.View(), "  \n\n\n\n", infoBox)

	//return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	s += m.footerView() //helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new test • q: exit\n"))
	return s
}

func (m UIinstance) headerView() string {
	title := "gobius" //titleStyle.Render("gobius")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m UIinstance) footerView() string {
	return helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new test • q: exit\n"))
	// info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	// line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	// return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

type LogViewWriter struct {
	App *tea.Program
}

func (tw *LogViewWriter) Write(p []byte) (n int, err error) {
	if tw.App == nil {
		return os.Stderr.Write(p)
	} else {
		tw.App.Send(string(p))
	}
	return len(p), nil
}

func InitialModel(validator IValidatorManager) tea.Model {

	return UIinstance{
		content:   "",
		ready:     false,
		viewport:  viewport.Model{},
		validator: validator,
	}

}

/*p := tea.NewProgram(
	model{content: string("")},
	tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
	//tea.WithMouseAllMotion(), // turn on mouse support so we can track the mouse wheel
)

go func() {
	for {
		logger.Warn().Msg("[red]Hello, World![white]")
		time.Sleep(250 * time.Millisecond)
	}
}()

if _, err := p.Run(); err != nil {
	fmt.Println("could not run program:", err)
	os.Exit(1)
}

*/
