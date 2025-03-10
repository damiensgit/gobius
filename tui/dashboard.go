package tui

import (
	"fmt"
	"gobius/metrics"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors matching btop
	subtle    = lipgloss.AdaptiveColor{Light: "#383838", Dark: "#2B2B2B"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#FD8F3F", Dark: "#FF9F4A"}
	critical  = lipgloss.AdaptiveColor{Light: "#FF5F5F", Dark: "#FF5F5F"}

	// Base styles
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	// Title style
	dashboardTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(subtle).
				Bold(true).
				MarginLeft(1).
				MarginRight(1)

	// GPU box style
	gpuBoxStyle = baseStyle.Copy().
			BorderForeground(highlight).
			Width(30)

	// Progress bar styles
	usageBarStyle = lipgloss.NewStyle().
			Foreground(special)

	memoryBarStyle = lipgloss.NewStyle().
			Foreground(warning)

	temperatureBarStyle = lipgloss.NewStyle().
				Foreground(critical)

	// Metrics box styles
	metricsBoxStyle = baseStyle.Copy().
			BorderForeground(highlight).
			Padding(1)

	// Value styles
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	// Table styles
	tableStyle = baseStyle.Copy().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1)
)

type GPUInfo struct {
	ID          int
	Temperature float64
	Usage       float64
	Memory      float64
	Status      string
	SolveTime   time.Duration
}

type ValidatorMetrics struct {
	SessionTime         string
	SolvedLastMinute    int64
	SolutionsLastMinute struct {
		Success int64
		Total   int64
		Rate    float64
	}
	AverageSolutionRate    float64
	AverageSolutionsPerMin float64
	AverageSolvesPerMin    float64
}

type FinancialMetrics struct {
	TokenIncomePerMin  float64
	TokenIncomePerHour float64
	TokenIncomePerDay  float64
	IncomePerMin       float64
	IncomePerHour      float64
	IncomePerDay       float64
	ProfitPerMin       float64
	ProfitPerHour      float64
	ProfitPerDay       float64
}

type model struct {
	width, height    int
	gpus             []metrics.GPUInfo
	validatorMetrics ValidatorMetrics
	financialMetrics FinancialMetrics
	progress         progress.Model
	table            table.Model
	logViewer        *LogViewer
	gpuViewer        *GPUViewer // Add this

	showingLogs      bool
	quitting         bool
	quitSignal       chan struct{}
	shutdownComplete bool
}

type updateMetricsMsg struct {
	validator ValidatorMetrics
	financial FinancialMetrics
}

type updateGPUMsg struct {
	gpus []metrics.GPUInfo
}

func NewDashboard() model {
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 30

	columns := []table.Column{
		{Title: "Metric", Width: 20},
		{Title: "Value", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(table.Styles{
		Header:   tableStyle.Bold(true),
		Selected: tableStyle.Foreground(highlight),
		Cell:     tableStyle,
	})

	return model{
		progress:         p,
		table:            t,
		logViewer:        NewLogViewer(),
		gpuViewer:        NewGPUViewer(),
		showingLogs:      false,
		quitting:         false,
		quitSignal:       make(chan struct{}),
		shutdownComplete: false,
		gpus:             []metrics.GPUInfo{
			// {ID: 0, Status: "Mining", SolveTime: 2 * time.Second},
			// {ID: 1, Status: "Mining", SolveTime: time.Duration(2.2 * float64(time.Second))},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate height for GPU viewer (adjust these numbers as needed)
		//gpuHeight := m.height/2 - 4 // Use top half of screen minus some padding
		//m.gpuViewer.SetSize(m.width-4, gpuHeight)

		if m.showingLogs {
			cmd = m.logViewer.Update(msg)
		} else {
			cmd = m.gpuViewer.Update(msg)
		}
	case tea.KeyMsg:
		if m.quitting {
			if m.shutdownComplete {
				return m, tea.Quit
			}
			return m, nil // Ignore all input while quitting
		}
		if m.showingLogs {
			switch msg.String() {
			case "q", "ctrl+c":
				m.quitting = true
				close(m.quitSignal)
				return m, nil
			case "esc":
				m.showingLogs = false
			default:
				cmd = m.logViewer.Update(msg)
			}
		} else {
			switch msg.String() {
			case "q", "ctrl+c":
				m.quitting = true
				close(m.quitSignal)
				return m, nil
			case "l":
				m.showingLogs = true
				m.logViewer.Update(tea.WindowSizeMsg{
					Width:  m.width,
					Height: m.height,
				})
				m.logViewer.RefreshContent()
			default:
				// Pass other key events to the GPU viewer when not showing logs
				cmd = m.gpuViewer.Update(msg)
			}
		}
	case tea.MouseMsg:
		if m.showingLogs {
			cmd = m.logViewer.Update(msg)
		} else {
			// Pass mouse events to the GPU viewer when not showing logs
			cmd = m.gpuViewer.Update(msg)
		}
	case updateMetricsMsg:
		m.validatorMetrics = msg.validator
		m.financialMetrics = msg.financial
	case updateGPUMsg:
		m.gpuViewer.UpdateGPUs(msg.gpus)
	case string:
		// Handle log messages
		if m.logViewer != nil {
			m.logViewer.AddLog(msg)
		}
	case tea.QuitMsg:
		m.Cleanup()
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.quitting {
		// Create a modal dialog style
		modalStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2).
			Width(50).
			Align(lipgloss.Center).
			MarginLeft((m.width - 50) / 2) // Center horizontally

		modal := modalStyle.Render("Shutting down - waiting for tasks to complete...\nPlease wait...")

		// Create vertical centering with newlines
		verticalPadding := (m.height - lipgloss.Height(modal)) / 2
		return strings.Repeat("\n", verticalPadding) + modal
	}

	if m.showingLogs {
		return m.logViewer.View()
	}

	var s strings.Builder

	// Title bar with help text
	title := fmt.Sprintf("Gobius Mining Dashboard - %s (press 'l' for logs, 'q' to quit)", m.validatorMetrics.SessionTime)
	s.WriteString(dashboardTitleStyle.Render(title))
	s.WriteString("\n\n")

	// GPU Section
	/*var gpuBoxes []string
	for _, gpu := range m.gpus {
		var box strings.Builder

		// GPU Title
		gpuTitle := fmt.Sprintf("GPU %d", gpu.ID)
		box.WriteString(lipgloss.NewStyle().Bold(true).Render(gpuTitle))
		box.WriteString("\n\n")

		// Status and solve time
		box.WriteString(fmt.Sprintf("Status: %s\n", gpu.Status))
		box.WriteString(fmt.Sprintf("Solve Time: %s", gpu.SolveTime))

		gpuBoxes = append(gpuBoxes, gpuBoxStyle.Render(box.String()))
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, gpuBoxes...))
	s.WriteString("\n\n")*/

	// GPU Section with scrolling
	s.WriteString(m.gpuViewer.View())
	s.WriteString("\n")

	// Create two columns for metrics
	var leftCol, rightCol strings.Builder

	// Validator Metrics Box
	validatorBox := "Validator Metrics\n\n"
	validatorBox += fmt.Sprintf("Solved Last Minute:     %s\n", valueStyle.Render(fmt.Sprintf("%d", m.validatorMetrics.SolvedLastMinute)))
	validatorBox += fmt.Sprintf("Solutions Success Rate: %s\n", valueStyle.Render(fmt.Sprintf("%.2f%%", m.validatorMetrics.SolutionsLastMinute.Rate*100)))
	validatorBox += fmt.Sprintf("Avg Solution Rate:      %s\n", valueStyle.Render(fmt.Sprintf("%.2f%%", m.validatorMetrics.AverageSolutionRate*100)))
	validatorBox += fmt.Sprintf("Avg Solutions/Min:      %s\n", valueStyle.Render(fmt.Sprintf("%.2f", m.validatorMetrics.AverageSolutionsPerMin)))
	validatorBox += fmt.Sprintf("Avg Solves/Min:         %s", valueStyle.Render(fmt.Sprintf("%.2f", m.validatorMetrics.AverageSolvesPerMin)))
	leftCol.WriteString(metricsBoxStyle.Render(validatorBox))

	// Financial Metrics Box
	financialBox := "Financial Metrics\n\n"
	financialBox += fmt.Sprintf("Token Income/Min:  %s AIUS\n", valueStyle.Render(fmt.Sprintf("%.4g", m.financialMetrics.TokenIncomePerMin)))
	financialBox += fmt.Sprintf("Token Income/Hour: %s AIUS\n", valueStyle.Render(fmt.Sprintf("%.4g", m.financialMetrics.TokenIncomePerHour)))
	financialBox += fmt.Sprintf("Token Income/Day:  %s AIUS\n\n", valueStyle.Render(fmt.Sprintf("%.4g", m.financialMetrics.TokenIncomePerDay)))
	financialBox += fmt.Sprintf("Income/Min:  %s\n", valueStyle.Render(fmt.Sprintf("$%.2f", m.financialMetrics.IncomePerMin)))
	financialBox += fmt.Sprintf("Income/Hour: %s\n", valueStyle.Render(fmt.Sprintf("$%.2f", m.financialMetrics.IncomePerHour)))
	financialBox += fmt.Sprintf("Income/Day:  %s\n\n", valueStyle.Render(fmt.Sprintf("$%.2f", m.financialMetrics.IncomePerDay)))
	financialBox += fmt.Sprintf("Profit/Min:  %s\n", valueStyle.Render(fmt.Sprintf("$%.2f", m.financialMetrics.ProfitPerMin)))
	financialBox += fmt.Sprintf("Profit/Hour: %s\n", valueStyle.Render(fmt.Sprintf("$%.2f", m.financialMetrics.ProfitPerHour)))
	financialBox += fmt.Sprintf("Profit/Day:  %s", valueStyle.Render(fmt.Sprintf("$%.2f", m.financialMetrics.ProfitPerDay)))
	rightCol.WriteString(metricsBoxStyle.Render(financialBox))

	// Join the columns
	s.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftCol.String(),
		"  ", // Add some spacing between columns
		rightCol.String(),
	))

	return s.String()
}

func (m model) UpdateValidatorMetrics(metrics ValidatorMetrics) tea.Cmd {
	return func() tea.Msg {
		return updateMetricsMsg{
			validator: metrics,
			financial: m.financialMetrics,
		}
	}
}

func (m model) UpdateFinancialMetrics(metrics FinancialMetrics) tea.Cmd {
	return func() tea.Msg {
		return updateMetricsMsg{
			validator: m.validatorMetrics,
			financial: metrics,
		}
	}
}

func (m model) UpdateGPUMetrics(gpus []metrics.GPUInfo) tea.Cmd {
	return func() tea.Msg {
		return updateGPUMsg{gpus: gpus}
	}
}

func (m *model) SendGPUMetrics(p *tea.Program, gpus []metrics.GPUInfo) {
	p.Send(updateGPUMsg{gpus: gpus}) // Directly sends the message to Bubble Tea
}

func (m *model) SetLogWriter(writer *LogViewWriter) {
	if m.logViewer != nil {
		m.logViewer.SetWriter(writer)
	}
}

func (m model) GetQuitSignal() chan struct{} {
	return m.quitSignal
}

// Add this method to handle final cleanup after tea.Quit
func (m *model) Cleanup() {
	m.shutdownComplete = true
	if m.logViewer != nil {
		m.logViewer.RestoreConsoleLogging()
	}
}
