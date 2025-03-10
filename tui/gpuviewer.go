package tui

import (
	"fmt"
	"math"
	"strings"

	"gobius/metrics"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMode int

const (
	DetailedView ViewMode = iota
	CompactView
	SummaryView
)

type GPUViewer struct {
	gpus             []metrics.GPUInfo
	width            int
	height           int
	scrollOffset     int
	viewMode         ViewMode
	compactThreshold int // Number of GPUs before switching to compact layout
}

// Common styles
var (
	baseCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1).
			MarginRight(1)

	gpuCardStyle = baseCardStyle.Copy().
			Width(20).
			MarginBottom(1)

	compactGPUCardStyle = baseCardStyle.Copy().
				Width(15).
				Height(2)

	gpuTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	gpuStatsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	scrollIndicator = lipgloss.NewStyle().
			Foreground(highlight).
			Align(lipgloss.Right)

	scrollUpIndicator   = scrollIndicator.Copy().SetString("↑")
	scrollDownIndicator = scrollIndicator.Copy().SetString("↓")

	gpuStatusStyles = map[string]lipgloss.Style{
		"Mining": lipgloss.NewStyle().Foreground(special),
		"Error":  lipgloss.NewStyle().Foreground(critical),
		"Idle":   lipgloss.NewStyle().Foreground(warning),
	}

	statusDots = map[string]string{
		"Mining": lipgloss.NewStyle().Foreground(special).SetString("●").String(),
		"Error":  lipgloss.NewStyle().Foreground(critical).SetString("●").String(),
		"Idle":   lipgloss.NewStyle().Foreground(subtle).SetString("●").String(),
	}
)

func NewGPUViewer() *GPUViewer {
	return &GPUViewer{
		gpus:             []metrics.GPUInfo{},
		scrollOffset:     0,
		viewMode:         DetailedView,
		compactThreshold: 16,
	}
}

// Helper function to render a row with scroll indicators
func (g *GPUViewer) renderRowWithScrollIndicators(content string, row, startRow, endRow, totalRows int) string {
	if g.scrollOffset > 0 && row == startRow {
		return lipgloss.JoinHorizontal(lipgloss.Top, content, scrollUpIndicator.Render())
	} else if endRow < totalRows && row == endRow-1 {
		return lipgloss.JoinHorizontal(lipgloss.Top, content, scrollDownIndicator.Render())
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, content, " ")
}

// Helper function to get view mode string
func (g *GPUViewer) getViewModeString() string {
	switch g.viewMode {
	case DetailedView:
		return "Detailed"
	case CompactView:
		return "Compact"
	case SummaryView:
		return "Summary"
	}
	return ""
}

func (g *GPUViewer) renderGPUs() string {
	if len(g.gpus) == 0 {
		return "No GPUs available"
	}

	// Count active GPUs
	active := 0
	for _, gpu := range g.gpus {
		if gpu.Status == "Mining" {
			active++
		}
	}

	// Render header
	var sb strings.Builder
	summary := fmt.Sprintf("GPUs: %d active / %d total | View: %s (press 'v' to change)",
		active, len(g.gpus), g.getViewModeString())
	sb.WriteString(lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(summary))
	sb.WriteString("\n")

	// Render view content
	switch g.viewMode {
	case SummaryView:
		sb.WriteString(g.renderSummaryView())
	case CompactView:
		sb.WriteString(g.renderCompactView())
	default:
		sb.WriteString(g.renderDetailedView())
	}

	return sb.String()
}

// Helper function to get layout metrics for all view modes
func (g *GPUViewer) getLayoutMetrics() (cardsPerRow, cardHeight, totalRows, visibleRows int) {
	switch g.viewMode {
	case DetailedView:
		cardHeight = 4
		cardsPerRow = (g.width - 4) / 22
	case CompactView:
		cardHeight = 3
		cardsPerRow = (g.width - 4) / 17
	case SummaryView:
		cardHeight = 1
		cardsPerRow = g.width - 6
	}

	if cardsPerRow < 1 {
		cardsPerRow = 1
	}

	totalRows = int(math.Ceil(float64(len(g.gpus)) / float64(cardsPerRow)))
	visibleRows = (g.height - 4) / cardHeight
	return
}

func (g *GPUViewer) getMaxScroll() int {
	_, _, totalRows, visibleRows := g.getLayoutMetrics()
	return max(0, totalRows-visibleRows)
}

func (g *GPUViewer) renderDetailedView() string {
	var sb strings.Builder

	cardsPerRow, cardHeight, totalRows, _ := g.getLayoutMetrics()
	rowsVisible := (g.height - 4) / cardHeight

	// Calculate visible range
	startRow := g.scrollOffset
	endRow := min(startRow+rowsVisible, totalRows)

	// Render visible rows
	for row := startRow; row < endRow; row++ {
		var currentRow []string
		startIdx := row * cardsPerRow
		endIdx := min((row+1)*cardsPerRow, len(g.gpus))

		for _, gpu := range g.gpus[startIdx:endIdx] {
			var cardContent strings.Builder

			cardContent.WriteString(gpuTitleStyle.Render(fmt.Sprintf("GPU %d", gpu.ID)))
			cardContent.WriteString("\n")

			statusStyle := gpuStatusStyles[gpu.Status]
			cardContent.WriteString(statusStyle.Render(gpu.Status))

			solveTimeStr := fmt.Sprintf("%0.1fs", gpu.SolveTime.Seconds())
			padding := 18 - len(gpu.Status) - len(solveTimeStr)
			if padding > 0 {
				cardContent.WriteString(strings.Repeat(" ", padding))
			}
			cardContent.WriteString(gpuStatsStyle.Render(solveTimeStr))

			currentRow = append(currentRow, gpuCardStyle.Render(cardContent.String()))
		}

		rowContent := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)

		// Add scroll indicators on the right if needed
		rowContent = g.renderRowWithScrollIndicators(rowContent, row, startRow, endRow, totalRows)

		sb.WriteString(rowContent)
		sb.WriteString("\n")
	}

	return sb.String()
}

func (g *GPUViewer) renderCompactView() string {
	var sb strings.Builder

	cardsPerRow, cardHeight, totalRows, _ := g.getLayoutMetrics()
	rowsVisible := (g.height - 4) / cardHeight

	// Calculate visible range
	startRow := g.scrollOffset
	endRow := min(startRow+rowsVisible, totalRows)

	// Render visible rows
	for row := startRow; row < endRow; row++ {
		var currentRow []string
		startIdx := row * cardsPerRow
		endIdx := min((row+1)*cardsPerRow, len(g.gpus))

		for _, gpu := range g.gpus[startIdx:endIdx] {
			var cardContent strings.Builder
			cardContent.WriteString(fmt.Sprintf("%d:%s %.1fs",
				gpu.ID,
				gpuStatusStyles[gpu.Status].Render(gpu.Status[:1]),
				gpu.SolveTime.Seconds()))

			currentRow = append(currentRow, compactGPUCardStyle.Render(cardContent.String()))
		}

		rowContent := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)

		// Add scroll indicators on the right if needed
		rowContent = g.renderRowWithScrollIndicators(rowContent, row, startRow, endRow, totalRows)

		sb.WriteString(rowContent)
		sb.WriteString("\n")
	}

	return sb.String()
}

func (g *GPUViewer) renderSummaryView() string {
	var sb strings.Builder

	cardsPerRow, _, totalRows, _ := g.getLayoutMetrics() // cardsPerRow is dotsPerRow in this case
	rowsVisible := g.height - 4

	startRow := g.scrollOffset
	endRow := min(startRow+rowsVisible, totalRows)

	for row := startRow; row < endRow; row++ {
		startIdx := row * cardsPerRow
		endIdx := min((row+1)*cardsPerRow, len(g.gpus))

		for _, gpu := range g.gpus[startIdx:endIdx] {
			sb.WriteString(statusDots[gpu.Status])
		}

		// Add scroll indicators on the right if needed
		sb.WriteString(g.renderRowWithScrollIndicators("", row, startRow, endRow, totalRows))

		sb.WriteString("\n")
	}

	return sb.String()
}

func (g *GPUViewer) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return g.handleKeyPress(msg)
	case tea.MouseMsg:
		return g.handleMouse(msg)
	case tea.WindowSizeMsg:
		g.SetSize(msg.Width-4, msg.Height/2-4)
	}
	return nil
}

func (g *GPUViewer) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "v":
		g.viewMode = (g.viewMode + 1) % 3
		g.scrollOffset = 0
	case "up", "k":
		if g.scrollOffset > 0 {
			g.scrollOffset--
		}
	case "down", "j":
		if g.scrollOffset < g.getMaxScroll() {
			g.scrollOffset++
		}
	case "home", "g":
		g.scrollOffset = 0
	case "end", "G":
		g.scrollOffset = g.getMaxScroll()
	}
	return nil
}

func (g *GPUViewer) handleMouse(msg tea.MouseMsg) tea.Cmd {
	switch msg.Type {
	case tea.MouseWheelUp:
		if g.scrollOffset > 0 {
			g.scrollOffset--
		}
	case tea.MouseWheelDown:
		if g.scrollOffset < g.getMaxScroll() {
			g.scrollOffset++
		}
	}
	return nil
}

func (g *GPUViewer) SetSize(width, height int) {
	g.width = width
	g.height = height
}

func (g *GPUViewer) UpdateGPUs(gpus []metrics.GPUInfo) {
	g.gpus = gpus

	// Only auto-switch view mode on initial load when empty
	if len(g.gpus) == 0 && len(gpus) > g.compactThreshold {
		g.viewMode = CompactView
	}
}

func (g *GPUViewer) View() string {
	return g.renderGPUs()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
