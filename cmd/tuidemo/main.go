package main

import (
	"fmt"
	"gobius/metrics"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

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

// Theme definitions
type Theme struct {
	Colors struct {
		Primary     tcell.Color
		Secondary   tcell.Color
		Accent      tcell.Color
		Background  tcell.Color
		Surface     tcell.Color
		Error       tcell.Color
		Warning     tcell.Color
		Success     tcell.Color
		Text        tcell.Color
		TextDim     tcell.Color
		Border      tcell.Color
		BorderFocus tcell.Color
		MenuFocus   tcell.Color
		MenuText    tcell.Color
		Title       tcell.Color
	}
	Borders struct {
		Horizontal  rune
		Vertical    rune
		TopLeft     rune
		TopRight    rune
		BottomLeft  rune
		BottomRight rune

		LeftT   rune
		RightT  rune
		TopT    rune
		BottomT rune
		Cross   rune

		HorizontalFocus  rune
		VerticalFocus    rune
		TopLeftFocus     rune
		TopRightFocus    rune
		BottomLeftFocus  rune
		BottomRightFocus rune
	}
	Symbols struct {
		BorderTopLeft     rune
		BorderTopRight    rune
		BorderBottomLeft  rune
		BorderBottomRight rune
		BorderHorizontal  rune
		BorderVertical    rune
		StatusMining      rune
		StatusError       rune
		StatusIdle        rune
		ScrollUp          rune
		ScrollDown        rune
		Bullet            rune
		ArrowUp           rune
		ArrowDown         rune
		ArrowFlat         rune
		MenuSeparator     rune
		TitleLeftDown     rune
		TitleRightDown    rune
		TitleLeft         rune
		TitleRight        rune
	}
	Styles struct {
		Title     tcell.Style
		Text      tcell.Style
		Highlight tcell.Style
		Dim       tcell.Style
		Error     tcell.Style
		Success   tcell.Style
		Warning   tcell.Style
	}
}

// Message types for state updates
type UpdateType int

const (
	UpdateGPUs UpdateType = iota
	UpdateLog
	UpdateValidatorMetrics
	UpdateFinancialMetrics
)

type ViewMode int

const (
	DetailedView ViewMode = iota
	CompactView
	SummaryView
	ViewDashboard
	ViewLogs
	ViewHelp
)

// Add back after UpdateType constants
type StateUpdate struct {
	Type    UpdateType
	Payload any
}

// First, let's create a struct to hold our utilization data
type GPUUtilization struct {
	Mining    int
	Error     int
	Idle      int
	Total     int
	Timestamp time.Time
}

// CustomBox is a box primitive with custom borders and title
type CustomBox struct {
	*tview.Box
	theme      *Theme
	title      string
	titleAlign int // Left: 0, Center: 1, Right: 2
}

func NewCustomBox(theme *Theme) *CustomBox {
	return &CustomBox{
		Box:        tview.NewBox(),
		theme:      theme,
		titleAlign: tview.AlignLeft,
	}
}

func (c *CustomBox) SetTitle(title string) *CustomBox {
	c.title = title
	return c
}

func (c *CustomBox) SetTitleAlign(align int) *CustomBox {
	c.titleAlign = align
	return c
}

func (c *CustomBox) Draw(screen tcell.Screen) {

	c.Box.DrawForSubclass(screen, c)

	// // Get box dimensions
	x, y, width, height := c.GetRect()

	// Draw border
	borderStyle := tcell.StyleDefault.
		Background(c.theme.Colors.Background).
		Foreground(c.theme.Colors.Border)

	// Draw corners
	screen.SetContent(x, y, c.theme.Symbols.BorderTopLeft, nil, borderStyle)
	screen.SetContent(x+width-1, y, c.theme.Symbols.BorderTopRight, nil, borderStyle)
	screen.SetContent(x, y+height-1, c.theme.Symbols.BorderBottomLeft, nil, borderStyle)
	screen.SetContent(x+width-1, y+height-1, c.theme.Symbols.BorderBottomRight, nil, borderStyle)

	// Draw horizontal borders
	for i := 1; i < width-1; i++ {
		screen.SetContent(x+i, y, c.theme.Symbols.BorderHorizontal, nil, borderStyle)
		screen.SetContent(x+i, y+height-1, c.theme.Symbols.BorderHorizontal, nil, borderStyle)
	}

	// Draw vertical borders
	for i := 1; i < height-1; i++ {
		screen.SetContent(x, y+i, c.theme.Symbols.BorderVertical, nil, borderStyle)
		screen.SetContent(x+width-1, y+i, c.theme.Symbols.BorderVertical, nil, borderStyle)
	}

	// Draw title if set
	if c.title != "" {
		titleStyle := tcell.StyleDefault.
			Background(c.theme.Colors.Background).
			Foreground(c.theme.Colors.Primary)

		// Calculate title position
		titleWidth := len(c.title) + 2
		var titleX int
		switch c.titleAlign {
		case tview.AlignCenter:
			titleX = x + (width-titleWidth)/2
		case tview.AlignRight:
			titleX = x + width - titleWidth - 3
		default: // AlignLeft
			titleX = x + 3
		}

		screen.SetContent(titleX, y, c.theme.Symbols.BorderTopRight, nil, borderStyle)

		// Draw title
		for i, r := range c.title {
			screen.SetContent(titleX+1+i, y, r, nil, titleStyle)
		}
		// Draw separator
		screen.SetContent(titleX+len(c.title)+1, y, c.theme.Symbols.BorderBottomLeft, nil, borderStyle)
	}

	// Draw inner conten
}

// Add these helper functions before the UtilizationMonitor struct definition

func round(x float64) float64 {
	return math.Floor(x + 0.5)
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

type UtilizationMonitor struct {
	*CustomBox
	samples    []GPUUtilization
	maxSamples int
	height     int // Height in rows
	width      int // Current width in columns
}

func NewUtilizationMonitor(theme *Theme, maxSamples, height int) *UtilizationMonitor {
	u := &UtilizationMonitor{
		CustomBox:  NewCustomBox(theme),
		samples:    make([]GPUUtilization, 0, maxSamples),
		maxSamples: maxSamples,
		height:     height,
		width:      0, // Will be set in Draw
	}
	u.SetTitle(" GPU Utilization ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetTitleColor(theme.Colors.Primary).
		SetBorderColor(theme.Colors.Border).
		SetBackgroundColor(theme.Colors.Background)

	return u
}

// Add a method to add new samples
func (u *UtilizationMonitor) AddSample(gpus []metrics.GPUInfo) {
	sample := GPUUtilization{
		Total:     len(gpus),
		Timestamp: time.Now(),
	}

	for _, gpu := range gpus {
		switch gpu.Status {
		case "Mining":
			sample.Mining++
		case "Error":
			sample.Error++
		case "Idle":
			sample.Idle++
		}
	}

	// Add sample and maintain history length based on current width
	u.samples = append(u.samples, sample)
	if len(u.samples) > u.maxSamples {
		u.samples = u.samples[1:]
	}
}

var graph_symbols = map[string][]rune{
	"braille_up": {
		' ', '⢀', '⢠', '⢰', '⢸',
		'⡀', '⣀', '⣠', '⣰', '⣸',
		'⡄', '⣄', '⣤', '⣴', '⣼',
		'⡆', '⣆', '⣦', '⣶', '⣾',
		'⡇', '⣇', '⣧', '⣷', '⣿',
	},
	"braille_down": {
		' ', '⠈', '⠘', '⠸', '⢸',
		'⠁', '⠉', '⠙', '⠹', '⢹',
		'⠃', '⠋', '⠛', '⠻', '⢻',
		'⠇', '⠏', '⠟', '⠿', '⢿',
		'⡇', '⡏', '⡟', '⡿', '⣿',
	},
	"block_up": {
		' ', '▗', '▗', '▐', '▐',
		'▖', '▄', '▄', '▟', '▟',
		'▖', '▄', '▄', '▟', '▟',
		'▌', '▙', '▙', '█', '█',
		'▌', '▙', '▙', '█', '█',
	},
	"block_down": {
		' ', '▝', '▝', '▐', '▐',
		'▘', '▀', '▀', '▜', '▜',
		'▘', '▀', '▀', '▜', '▜',
		'▌', '▛', '▛', '█', '█',
		'▌', '▛', '▛', '█', '█',
	},
	"tty_up": {
		' ', '░', '░', '▒', '▒',
		'░', '░', '▒', '▒', '█',
		'░', '▒', '▒', '▒', '█',
		'▒', '▒', '▒', '█', '█',
		'▒', '█', '█', '█', '█',
	},
	"tty_down": {
		' ', '░', '░', '▒', '▒',
		'░', '░', '▒', '▒', '█',
		'░', '▒', '▒', '▒', '█',
		'▒', '▒', '▒', '█', '█',
		'▒', '█', '█', '█', '█',
	},
}

// Add this color gradient array before the Draw method
var utilizationColors = []tcell.Color{
	tcell.NewRGBColor(255, 47, 47),  // Red
	tcell.NewRGBColor(255, 91, 47),  // Red-Orange
	tcell.NewRGBColor(255, 134, 47), // Orange
	tcell.NewRGBColor(255, 177, 47), // Orange-Yellow
	tcell.NewRGBColor(255, 221, 47), // Yellow
	tcell.NewRGBColor(221, 255, 47), // Yellow-Green
	tcell.NewRGBColor(177, 255, 47), // Light Green
	tcell.NewRGBColor(134, 255, 47), // Green
	tcell.NewRGBColor(91, 255, 47),  // Bright Green
	tcell.NewRGBColor(47, 255, 47),  // Pure Green
}

func (u *UtilizationMonitor) Draw(screen tcell.Screen) {
	// First draw the custom box (border and title)
	u.CustomBox.Draw(screen)

	// Get the inner area for drawing the chart
	x, y, width, height := u.GetInnerRect()

	// Check if width changed
	if width != u.width {
		u.width = width
		// Each braille character represents 2 samples
		newMaxSamples := width * 2
		// Adjust maxSamples
		u.maxSamples = newMaxSamples
		// Trim samples if we have too many
		if len(u.samples) > newMaxSamples {
			u.samples = u.samples[len(u.samples)-newMaxSamples:]
		}
	}

	// Calculate how many samples we can show based on width
	samplesToShow := min(len(u.samples), width*2)
	if samplesToShow == 0 {
		return
	}

	chartPatterns := graph_symbols["braille_up"]

	// Draw the utilization bars
	for i := 0; i < (samplesToShow+1)/2; i++ {
		// Calculate the indices for the two samples this character will represent
		sampleIdx1 := len(u.samples) - samplesToShow + (i * 2)
		sampleIdx2 := sampleIdx1 + 1

		// Get the samples
		sample1 := u.samples[sampleIdx1]
		var sample2 GPUUtilization
		if sampleIdx2 < len(u.samples) {
			sample2 = u.samples[sampleIdx2]
		}

		// Calculate utilization percentages (0-100)
		util1 := int64(float64(sample1.Mining) / float64(sample1.Total) * 100)
		util2 := int64(0)
		if sampleIdx2 < len(u.samples) {
			util2 = int64(float64(sample2.Mining) / float64(sample2.Total) * 100)
		}

		// For each row in the box
		for row := 0; row < height; row++ {
			var color tcell.Color

			if height == 1 {
				// Single row case - color based on average utilization
				avgUtil := float64(util1) / 100.0
				if sampleIdx2 < len(u.samples) {
					avgUtil = (float64(util1) + float64(util2)) / 200.0
				}
				colorIndex := int(avgUtil * float64(len(utilizationColors)-1))
				if colorIndex >= len(utilizationColors) {
					colorIndex = len(utilizationColors) - 1
				}
				color = utilizationColors[colorIndex]
			} else {
				// Multi-row case - color based on row position
				// Calculate which color this row should use based on its position from bottom to top
				rowFromBottom := height - 1 - row
				colorIndex := int(float64(rowFromBottom) * float64(len(utilizationColors)-1) / float64(height-1))
				if colorIndex >= len(utilizationColors) {
					colorIndex = len(utilizationColors) - 1
				}
				color = utilizationColors[colorIndex]
			}

			// Calculate the current row's value range (0-100)
			curHigh := int64(100)
			curLow := int64(0)
			if height > 1 {
				curHigh = int64(float64(height-row) * 100.0 / float64(height))
				curLow = int64(float64(height-(row+1)) * 100.0 / float64(height))
			}

			// Calculate braille pattern for both values
			var result [2]int
			const mod float64 = 0.1 // Adjustment factor from btop

			// Process first value (left column)
			if util1 >= curHigh {
				result[0] = 4
			} else if util1 <= curLow {
				result[0] = 0
			} else {
				val := float64(util1-curLow)*4.0/float64(curHigh-curLow) + mod
				result[0] = clamp(int(round(val)), 0, 4)
			}

			// Process second value (right column)
			if util2 >= curHigh {
				result[1] = 4
			} else if util2 <= curLow {
				result[1] = 0
			} else {
				val := float64(util2-curLow)*4.0/float64(curHigh-curLow) + mod
				result[1] = clamp(int(round(val)), 0, 4)
			}

			// Convert the two 0-4 values into braille pattern index using btop's braille_up pattern
			patternIndex := result[0]*5 + result[1]

			// Draw the braille character
			screen.SetContent(
				x+i,
				y+row,
				chartPatterns[patternIndex],
				nil,
				tcell.StyleDefault.Foreground(color),
			)
		}
	}
}

// Dashboard represents the main TUI application
type Dashboard struct {
	app                *tview.Application
	layout             *tview.Flex
	titleBar           *tview.Flex     // Add title bar
	leftTitle          *tview.TextView // Left title component
	rightTitle         *tview.TextView // Right title component
	menu               *tview.Flex     // Changed from *tview.TextView to *tview.Flex
	leftMenu           *tview.TextView // Add left menu component
	rightMenu          *tview.TextView // Add right menu component
	contentArea        *tview.Flex
	gpuViewer          *GPUViewer
	logViewer          *LogViewer
	metricsView        *MetricsView
	helpView           *HelpViewer
	theme              *Theme
	updates            chan StateUpdate
	mu                 sync.RWMutex
	validatorMetrics   ValidatorMetrics
	financialMetrics   FinancialMetrics
	currentView        ViewMode
	utilizationMonitor *UtilizationMonitor
	logger             zerolog.Logger
}

// Updated GPUViewer struct
type GPUViewer struct {
	*CustomBox
	app              *tview.Application
	theme            *Theme
	gpus             []metrics.GPUInfo
	viewMode         ViewMode
	compactThreshold int
	table            *tview.Table
	container        *tview.Flex     // New container to hold table and details panel
	detailsPanel     *tview.TextView // Panel for GPU details
	hasFocus         bool
	lastClickTime    time.Time
}

// Updated NewGPUViewer to initialize container and update InputCapture for spacebar
func NewGPUViewer(app *tview.Application, theme *Theme) *GPUViewer {
	g := &GPUViewer{
		CustomBox:        NewCustomBox(theme),
		app:              app,
		theme:            theme,
		viewMode:         DetailedView,
		compactThreshold: 16,
		table:            tview.NewTable(),
	}

	// Configure table defaults
	g.table.
		SetSelectable(true, false).
		SetBackgroundColor(theme.Colors.Background).
		SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.MouseLeftClick {
				x, y := event.Position()
				tableX, tableY, _, _ := g.table.GetRect()
				row := y - tableY
				col := x - tableX
				if g.viewMode == DetailedView && row == 0 {
					// Don't select header row
					return action, nil
				}
				g.table.Select(row, col)
			}
			return action, event
		})

	g.SetTitle(" GPUs ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetTitleColor(theme.Colors.Primary).
		SetBorderColor(theme.Colors.Border).
		SetBackgroundColor(theme.Colors.Background)

	// Initialize container as a Flex with vertical direction, add table as first item
	g.container = tview.NewFlex().SetDirection(tview.FlexRow)
	g.container.AddItem(g.table, 0, 1, true)

	// Update table InputCapture to also handle spacebar to toggle details
	g.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'v' || event.Rune() == 'V' {
				g.viewMode = (g.viewMode + 1) % 3
				g.UpdateGPUs(g.gpus)
				return nil
			} else if event.Rune() == ' ' {
				g.toggleDetails()
				return nil
			}
		}
		return event
	})

	g.table.SetSelectionChangedFunc(func(row, col int) {
		if g.detailsPanel != nil {
			g.updateDetailsPanel(row, col)
		}
	})

	return g
}

// Implement MouseHandler for GPUViewer
func (g *GPUViewer) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		// Get the table's mouse handler to delegate unhandled events
		if tableHandler := g.table.MouseHandler(); tableHandler != nil {
			// // Get the inner area coordinates
			// x, y, width, height := g.GetInnerRect()
			// mouseX, mouseY := event.Position()

			// // Check if click is within our bounds
			// if mouseX < x || mouseX >= x+width || mouseY < y || mouseY >= y+height {
			// 	return false, nil
			// }

			// if action == tview.MouseLeftClick {
			// 	setFocus(g)

			// 	// Calculate relative position within the table
			// 	relativeY := mouseY - y
			// 	relativeX := mouseX - x

			// 	// For DetailedView and CompactView, select the row
			// 	if g.viewMode != SummaryView {
			// 		if g.viewMode == DetailedView && relativeY == 0 {
			// 			// Don't select header row
			// 			return true, nil
			// 		}
			// 		g.table.Select(relativeY, 0)
			// 		return true, nil
			// 	}

			// 	// For SummaryView, select the cell
			// 	cellWidth := 3
			// 	col := relativeX / cellWidth
			// 	g.table.Select(relativeY, col)
			// 	return true, nil
			// }

			// Delegate all other mouse actions (like scrolling) to the table's handler
			return tableHandler(action, event, setFocus)
		}
		return false, nil
	}
}

// Helper function to get GPU index from current selection
func (g *GPUViewer) getGPUIndexFromSelection() int {
	row, col := g.table.GetSelection()
	if g.viewMode == DetailedView {
		return row - 1 // Account for header row
	} else if g.viewMode == SummaryView {
		_, _, width, _ := g.table.GetRect()
		cols := width / 3 // 3 is cell width in summary view
		if cols < 1 {
			cols = 1
		}
		return row*cols + col
	}
	return row
}

// Implement Focus method from tview.Primitive
func (g *GPUViewer) Focus(delegate func(p tview.Primitive)) {
	delegate(g.table)
	g.hasFocus = true
}

// Implement HasFocus method from tview.Primitive
func (g *GPUViewer) HasFocus() bool {
	return g.hasFocus || g.table.HasFocus()
}

// Implement Blur method from tview.Primitive
func (g *GPUViewer) Blur() {
	g.hasFocus = false
	g.table.Blur()
}

// Implement InputHandler to delegate keyboard events to the underlying table and provide basic arrow key scrolling
func (g *GPUViewer) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	if h := g.table.InputHandler(); h != nil {
		return h
	}
	return nil
}

// Update GPUViewer.Draw to render the container
func (g *GPUViewer) Draw(screen tcell.Screen) {
	g.CustomBox.Draw(screen)

	// Get the box's inner dimensions
	x, y, width, height := g.GetInnerRect()

	if g.detailsPanel != nil {
		outerX, _, outerWidth, outerHeight := g.GetRect()

		// Draw border
		borderStyle := tcell.StyleDefault.
			Background(g.theme.Colors.Background).
			Foreground(g.theme.Colors.Border)

		// Draw horizontal borders
		for i := outerX + 1; i < outerWidth-1; i++ {
			screen.SetContent(i, outerHeight-4, g.theme.Symbols.BorderHorizontal, nil, borderStyle)
		}
		screen.SetContent(outerX, outerHeight-4, g.theme.Symbols.TitleRight, nil, borderStyle)
		screen.SetContent(outerWidth-1, outerHeight-4, g.theme.Symbols.TitleLeft, nil, borderStyle)

		g.detailsPanel.SetRect(x, height-1, width, 3)
		g.detailsPanel.Draw(screen)

		height -= 4
	}

	// Set container dimensions and draw it
	g.container.SetRect(x, y, width, height)
	g.container.Draw(screen)
}

// New method to toggle splitting GPU view to show/hide details panel
func (g *GPUViewer) toggleDetails() {
	if g.detailsPanel == nil {
		// Create a simple details text view without a separate box (no border/title)
		details := tview.NewTextView().SetDynamicColors(true)
		// Set the details panel
		g.detailsPanel = details
		// Immediately update the details panel with the current selection
		row, col := g.table.GetSelection()
		g.updateDetailsPanel(row, col)
	} else {
		// Remove details panel
		g.detailsPanel = nil
	}
}

func (g *GPUViewer) updateDetailsPanel(row, col int) {
	if g.detailsPanel != nil {
		var index int
		if g.viewMode == DetailedView {
			index = row - 1
			if index < 0 {
				index = 0
			}
		} else if g.viewMode == SummaryView {
			_, _, fullWidth, _ := g.table.GetRect()
			cellWidth := 3
			cols := fullWidth / cellWidth
			if cols < 1 {
				cols = 1
			}
			index = row*cols + col
		} else {
			index = row
		}
		if index >= 0 && index < len(g.gpus) {
			gpu := g.gpus[index]
			detailText := fmt.Sprintf("[gray]────────────────────────────[-]\nGPU ID: %d\nStatus: %s\nSolve Time: %.1fs\n", gpu.ID, gpu.Status, gpu.SolveTime.Seconds())
			g.detailsPanel.SetText(detailText)
		} else {
			g.detailsPanel.SetText("[gray]────────────────────────────[-]\nNo GPU selected")
		}
	}
}

func (g *GPUViewer) renderDetailedView() {
	g.table.Clear()

	// Configure for row-only selection
	g.table.SetSelectable(true, false)

	// Set up header
	headerStyle := g.theme.Styles.Title
	g.table.SetCell(0, 0, tview.NewTableCell("GPU").SetStyle(headerStyle).SetSelectable(false))
	g.table.SetCell(0, 1, tview.NewTableCell("Status").SetStyle(headerStyle).SetSelectable(false))
	g.table.SetCell(0, 2, tview.NewTableCell("Solve Time").SetStyle(headerStyle).SetSelectable(false))
	g.table.SetCell(0, 3, tview.NewTableCell("History").SetStyle(headerStyle).SetSelectable(false))

	// Fix the header row
	g.table.SetFixed(1, 0)

	sparklineWidth := 10

	// Render all GPUs (starting from row 1 because row 0 is header)
	for i := 0; i < len(g.gpus); i++ {
		gpu := &g.gpus[i]
		row := i + 1

		// GPU ID
		idCell := tview.NewTableCell(fmt.Sprintf("%d", gpu.ID)).
			SetStyle(g.theme.Styles.Text).
			SetExpansion(1).
			SetAlign(tview.AlignLeft)

		// Status with symbol
		statusSymbol := g.theme.Symbols.StatusMining
		statusStyle := g.theme.Styles.Success

		switch gpu.Status {
		case "Error":
			statusSymbol = g.theme.Symbols.StatusError
			statusStyle = g.theme.Styles.Error
		case "Idle":
			statusSymbol = g.theme.Symbols.StatusIdle
			statusStyle = g.theme.Styles.Warning
		}

		statusCell := tview.NewTableCell(fmt.Sprintf("%c %s", statusSymbol, gpu.Status)).
			SetStyle(statusStyle).
			SetExpansion(1).
			SetAlign(tview.AlignLeft)

		// Solve time
		timeCell := tview.NewTableCell(fmt.Sprintf("%.1fs", gpu.SolveTime.Seconds())).
			SetStyle(g.theme.Styles.Dim).
			SetExpansion(1).
			SetAlign(tview.AlignLeft)

		sparkline := generateSparkline(gpu.SolveTimes, sparklineWidth)
		sparklineCell := tview.NewTableCell(sparkline).
			SetStyle(g.theme.Styles.Text).
			SetExpansion(2).
			SetAlign(tview.AlignLeft)

		g.table.SetCell(row, 0, idCell)
		g.table.SetCell(row, 1, statusCell)
		g.table.SetCell(row, 2, timeCell)
		g.table.SetCell(row, 3, sparklineCell)
	}
}

func (g *GPUViewer) renderCompactView() {
	g.table.Clear()

	// Make table scrollable
	g.table.SetSelectable(true, false).
		SetEvaluateAllRows(true)

	cardsPerRow := 4
	cellWidth := 12 // Fixed width for each cell

	for i := 0; i < len(g.gpus); i++ {
		gpu := &g.gpus[i]
		row := i / cardsPerRow
		col := i % cardsPerRow

		statusStyle := g.theme.Styles.Success
		switch gpu.Status {
		case "Error":
			statusStyle = g.theme.Styles.Error
		case "Idle":
			statusStyle = g.theme.Styles.Warning
		}

		// Create a fixed-width cell with status-colored background
		text := fmt.Sprintf("%d:%s %.1fs", gpu.ID, gpu.Status[:1], gpu.SolveTime.Seconds())
		cell := tview.NewTableCell(text).
			SetStyle(statusStyle).
			SetSelectable(true).
			SetExpansion(1).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(cellWidth)

		g.table.SetCell(row, col, cell)
	}
}

func (g *GPUViewer) renderSummaryView() {
	g.table.Clear()

	// Enable both row and column selection for cell-based selection
	g.table.SetSelectable(true, true)

	// Set selected function to show details on Enter/click
	g.table.SetSelectedFunc(func(row, col int) {
		index := g.getGPUIndexFromSelection()
		if index < len(g.gpus) {
			g.showGpuDetails(index)
		}
	})

	// Determine cell width and number of columns per row
	cellWidth := 3 // width for each gpu cell; adjust as needed
	_, _, fullWidth, _ := g.table.GetRect()
	colsPerRow := fullWidth / cellWidth
	if colsPerRow < 1 {
		colsPerRow = 1
	}

	for i, gpu := range g.gpus {
		row := i / colsPerRow
		col := i % colsPerRow
		var symbol string
		var style tcell.Style
		switch gpu.Status {
		case "Mining":
			// Animate mining status using rotating frames
			frames := []string{"◐", "◓", "◑", "◒"}
			frameIndex := (int(time.Now().UnixNano() / 200000000)) % len(frames) // change frame every 200ms
			symbol = frames[frameIndex]
			style = g.theme.Styles.Success
		case "Error":
			symbol = "✖"
			style = g.theme.Styles.Error
		case "Idle":
			symbol = "○"
			style = g.theme.Styles.Warning
		default:
			symbol = "●"
			style = g.theme.Styles.Text
		}

		cell := tview.NewTableCell(symbol).
			SetStyle(style).
			SetAlign(tview.AlignCenter).
			SetSelectable(true).
			SetExpansion(1).
			SetMaxWidth(cellWidth)

		g.table.SetCell(row, col, cell)
	}
}

func (g *GPUViewer) showGpuDetails(index int) {
	gpu := g.gpus[index]
	detail := fmt.Sprintf("GPU ID: %d\nStatus: %s\nSolve Time: %.1fs\n", gpu.ID, gpu.Status, gpu.SolveTime.Seconds())
	modal := tview.NewModal().
		SetText(detail).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.app.SetRoot(g.table, true).SetFocus(g.table)
		})
	g.app.SetRoot(modal, false).SetFocus(modal)
}

// HelpViewer component
type HelpViewer struct {
	*CustomTextView
	theme *Theme
}

func NewHelpViewer(theme *Theme) *HelpViewer {

	helpView := NewCustomTextView(theme)
	helpView.
		SetTitle(" Help ").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitleColor(theme.Colors.Primary).
		SetBorderColor(theme.Colors.Border)

	return &HelpViewer{
		CustomTextView: helpView,
		theme:          theme,
	}
}

// LogViewer component
type LogViewer struct {
	*CustomTextView
	theme    *Theme
	maxLines int
}

func (w *LogViewer) Write(p []byte) (n int, err error) {
	// Ensure logs appear immediately in the TextView
	//fmt.Fprintf(w.TextView, "%s", p)
	return w.TextView.Write(p)
	//return len(p), nil
	//return tview.ANSIWriter(w.TextView).Write(p)
}

func NewLogsViewer(theme *Theme) *LogViewer {

	logsView := NewCustomTextView(theme)
	logsView.
		SetTitle(" Logs ").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitleColor(theme.Colors.Primary).
		SetBorderColor(theme.Colors.Border).
		SetScrollable(true).ScrollToEnd()

	return &LogViewer{
		CustomTextView: logsView,
		theme:          theme,
		maxLines:       1000,
	}
}

// MetricsView component
type MetricsView struct {
	*CustomTextView
	theme            *Theme
	validatorMetrics ValidatorMetrics
	financialMetrics FinancialMetrics
}

func NewMetricsView(theme *Theme) *MetricsView {
	return &MetricsView{
		CustomTextView: NewCustomTextView(theme),
		theme:          theme,
	}
}

// CustomTextView is a TextView with custom borders and title
type CustomTextView struct {
	*tview.TextView
	customBox *CustomBox
	theme     *Theme
}

// NewCustomTextView creates a new TextView with our CustomBox
func NewCustomTextView(theme *Theme) *CustomTextView {
	textView := tview.NewTextView()
	customBox := NewCustomBox(theme)

	textView.Box = customBox.Box

	textView.Box.SetBorder(false)

	// Set common properties
	textView.SetDynamicColors(true)
	textView.SetScrollable(true)
	textView.SetWrap(true)

	// Set theme colors
	textView.SetBackgroundColor(theme.Colors.Background)
	textView.SetTextColor(theme.Colors.Text)

	return &CustomTextView{
		TextView:  textView,
		customBox: customBox,
		theme:     theme,
	}
}

// Draw overrides TextView's Draw method to use the custom border
func (ctv *CustomTextView) Draw(screen tcell.Screen) {
	// Update the CustomBox with the TextView's dimensions
	ctv.customBox.Box.SetRect(ctv.TextView.GetRect())

	// Draw the custom box with its fancy borders
	ctv.customBox.Draw(screen)

	// Get the inner area where text should be drawn
	x, y, width, height := ctv.customBox.GetInnerRect()

	// Tell the TextView to draw only in the inner area
	ctv.TextView.SetRect(x+1, y+1, width-2, height-2)

	// // Draw the TextView content (without its borders)
	ctv.TextView.Box.SetBorder(false)
	//ctv.TextView.Box.SetBorderPadding(1, 1, 1, 1)

	ctv.TextView.Draw(screen)
}

// GetInnerRect ensures text is drawn within custom borders
func (ctv *CustomTextView) GetInnerRect() (int, int, int, int) {
	return ctv.customBox.GetInnerRect()
}

// GetRect returns the outer rectangle dimensions
func (ctv *CustomTextView) GetRect() (int, int, int, int) {
	return ctv.customBox.GetRect()
}

// // SetRect sets the outer rectangle dimensions
// func (ctv *CustomTextView) SetRect(x, y, width, height int) *CustomTextView {
// 	ctv.customBox.Box.SetRect(x, y, width, height)
// 	return ctv
// }

// InputHandler returns the handler for keyboard events
func (ctv *CustomTextView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ctv.TextView.InputHandler()
}

// Focus is called when this primitive receives focus
func (ctv *CustomTextView) Focus(delegate func(p tview.Primitive)) {
	ctv.TextView.Focus(delegate)
}

// HasFocus returns whether or not this primitive has focus
func (ctv *CustomTextView) HasFocus() bool {
	return ctv.TextView.HasFocus()
}

// Blur is called when this primitive loses focus
func (ctv *CustomTextView) Blur() {
	ctv.TextView.Blur()
}

// MouseHandler returns the handler for mouse events
func (ctv *CustomTextView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return ctv.TextView.MouseHandler()
}

// SetTitle sets the title of the custom box
func (ctv *CustomTextView) SetTitle(title string) *CustomTextView {
	ctv.customBox.SetTitle(title)
	//ctv.TextView.SetTitle(title) // Also set on TextView for consistency
	return ctv
}

// SetTitleAlign sets the title alignment of the custom box
func (ctv *CustomTextView) SetTitleAlign(align int) *CustomTextView {
	ctv.customBox.SetTitleAlign(align)
	//ctv.TextView.SetTitleAlign(align) // Also set on TextView for consistency
	return ctv
}

// SetBorder sets the border of both the custom box and the text view
func (ctv *CustomTextView) SetBorder(show bool) *CustomTextView {
	ctv.customBox.SetBorder(show)
	return ctv
}

// SetBorderColor sets the border color of the custom box
func (ctv *CustomTextView) SetBorderColor(color tcell.Color) *CustomTextView {
	ctv.customBox.SetBorderColor(color)
	return ctv
}

// SetTitleColor sets the title color of the custom box
func (ctv *CustomTextView) SetTitleColor(color tcell.Color) *CustomTextView {
	ctv.customBox.SetTitleColor(color)
	return ctv
}

// SetTextAlign sets the text alignment of the text view
func (ctv *CustomTextView) SetTextAlign(align int) *CustomTextView {
	ctv.TextView.SetTextAlign(align)
	return ctv
}

// SetText sets the text of the text view
func (ctv *CustomTextView) SetText(text string) *CustomTextView {
	ctv.TextView.SetText(text)
	return ctv
}

// Initialize the modern theme
var ModernTheme = &Theme{
	Colors: struct {
		Primary     tcell.Color
		Secondary   tcell.Color
		Accent      tcell.Color
		Background  tcell.Color
		Surface     tcell.Color
		Error       tcell.Color
		Warning     tcell.Color
		Success     tcell.Color
		Text        tcell.Color
		TextDim     tcell.Color
		Border      tcell.Color
		BorderFocus tcell.Color
		MenuFocus   tcell.Color
		MenuText    tcell.Color
		Title       tcell.Color
	}{
		Primary:     tcell.NewRGBColor(74, 40, 255),   // Arbius blue #4a 28 ff
		Secondary:   tcell.NewRGBColor(138, 180, 248), // Lighter blue
		Accent:      tcell.NewRGBColor(412, 98, 255),  // Bright blue
		Background:  tcell.NewRGBColor(0, 0, 0),       // Pure black
		Surface:     tcell.NewRGBColor(18, 18, 18),    // Very dark gray
		Error:       tcell.NewRGBColor(234, 67, 53),   // Google-style red
		Warning:     tcell.NewRGBColor(251, 188, 4),   // Amber warning
		Success:     tcell.NewRGBColor(66, 133, 244),  // Blue for success
		Text:        tcell.NewRGBColor(242, 242, 242), // Almost white
		TextDim:     tcell.NewRGBColor(128, 128, 128), // Medium gray
		Border:      tcell.NewRGBColor(130, 130, 130), // Dark border
		BorderFocus: tcell.NewRGBColor(66, 133, 244),  // Blue for focus
		MenuFocus:   tcell.NewRGBColor(66, 133, 244),  // Blue for menu focus
		MenuText:    tcell.NewRGBColor(242, 242, 242), // Light gray for menu text
		Title:       tcell.NewRGBColor(74, 40, 255),   // Arbius blue for titles
	},
	Borders: struct {
		Horizontal  rune
		Vertical    rune
		TopLeft     rune
		TopRight    rune
		BottomLeft  rune
		BottomRight rune

		LeftT   rune
		RightT  rune
		TopT    rune
		BottomT rune
		Cross   rune

		HorizontalFocus  rune
		VerticalFocus    rune
		TopLeftFocus     rune
		TopRightFocus    rune
		BottomLeftFocus  rune
		BottomRightFocus rune
	}{
		Horizontal:  tview.BoxDrawingsLightHorizontal,
		Vertical:    tview.BoxDrawingsLightVertical,
		TopLeft:     tview.BoxDrawingsLightArcDownAndRight,
		TopRight:    tview.BoxDrawingsLightArcDownAndLeft,
		BottomLeft:  tview.BoxDrawingsLightArcUpAndRight,
		BottomRight: tview.BoxDrawingsLightArcUpAndLeft,

		LeftT:   tview.BoxDrawingsLightVerticalAndRight,
		RightT:  tview.BoxDrawingsLightVerticalAndLeft,
		TopT:    tview.BoxDrawingsLightDownAndHorizontal,
		BottomT: tview.BoxDrawingsLightUpAndHorizontal,
		Cross:   tview.BoxDrawingsLightVerticalAndHorizontal,

		HorizontalFocus:  tview.BoxDrawingsHeavyHorizontal,
		VerticalFocus:    tview.BoxDrawingsHeavyVertical,
		TopLeftFocus:     tview.BoxDrawingsHeavyDownAndRight,
		TopRightFocus:    tview.BoxDrawingsHeavyDownAndLeft,
		BottomLeftFocus:  tview.BoxDrawingsHeavyUpAndRight,
		BottomRightFocus: tview.BoxDrawingsHeavyUpAndLeft,
	},
	Symbols: struct {
		BorderTopLeft     rune
		BorderTopRight    rune
		BorderBottomLeft  rune
		BorderBottomRight rune
		BorderHorizontal  rune
		BorderVertical    rune
		StatusMining      rune
		StatusError       rune
		StatusIdle        rune
		ScrollUp          rune
		ScrollDown        rune
		Bullet            rune
		ArrowUp           rune
		ArrowDown         rune
		ArrowFlat         rune
		MenuSeparator     rune
		TitleLeftDown     rune
		TitleRightDown    rune
		TitleLeft         rune
		TitleRight        rune
	}{
		BorderTopLeft:     '╭',
		BorderTopRight:    '╮',
		BorderBottomLeft:  '╰',
		BorderBottomRight: '╯',
		BorderHorizontal:  '─',
		BorderVertical:    '│',
		StatusMining:      '◉',
		StatusError:       '✖',
		StatusIdle:        '○',
		ScrollUp:          '▲',
		ScrollDown:        '▼',
		Bullet:            '•',
		ArrowUp:           '↑',
		ArrowDown:         '↓',
		ArrowFlat:         '→',
		MenuSeparator:     '│',
		TitleLeftDown:     '┘',
		TitleRightDown:    '└',
		TitleLeft:         '┤',
		TitleRight:        '├',
	},
}

func (t *Theme) InitStyles() {

	tview.Borders = t.Borders

	t.Styles.Title = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.Primary).
		Bold(true)

	t.Styles.Text = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.Text)

	t.Styles.Highlight = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.Accent).
		Bold(true)

	t.Styles.Dim = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.TextDim)

	t.Styles.Error = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.Error)

	t.Styles.Success = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.Success)

	t.Styles.Warning = tcell.StyleDefault.
		Background(t.Colors.Background).
		Foreground(t.Colors.Warning)
}

// Add constants for minimum screen dimensions
const (
	MinScreenWidth  = 80
	MinScreenHeight = 24
)

func NewDashboard() *Dashboard {
	theme := ModernTheme
	theme.InitStyles()

	d := &Dashboard{
		app:     tview.NewApplication().EnableMouse(true),
		theme:   theme,
		updates: make(chan StateUpdate, 100),
	}

	d.initializeComponents()
	d.setupLayout()
	d.startUpdateHandler()

	return d
}

func (d *Dashboard) GetLogger() zerolog.Logger {
	return d.logger
}

func (d *Dashboard) initializeComponents() {
	// Initialize Title Bar components
	d.leftTitle = tview.NewTextView()
	d.leftTitle.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[#ffffff]GOBIUS[-]").
		SetBackgroundColor(d.theme.Colors.Primary)

	d.rightTitle = tview.NewTextView()
	d.rightTitle.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetText("[#ffffff]v1.0.0[-]").
		SetTextStyle(tcell.StyleDefault.Background(d.theme.Colors.Primary).Foreground(d.theme.Colors.Title)).
		SetBackgroundColor(d.theme.Colors.Primary)

	// Create a Flex container for the title bar
	d.titleBar = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.leftTitle, 0, 1, false).  // Expandable left title
		AddItem(d.rightTitle, 10, 0, false) // Fixed width right title
	d.titleBar.SetBackgroundColor(d.theme.Colors.Primary)

	// Initialize Menu components
	d.leftMenu = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	d.leftMenu.SetBackgroundColor(d.theme.Colors.Secondary)

	d.rightMenu = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	d.rightMenu.SetBackgroundColor(d.theme.Colors.Secondary)

	// Create a Flex container for the menu
	d.menu = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(d.leftMenu, 0, 1, false).  // Expandable left menu
		AddItem(d.rightMenu, 10, 0, false) // Fixed width right menu
	d.menu.SetBackgroundColor(d.theme.Colors.Secondary)

	// Initialize Help View
	d.helpView = NewHelpViewer(d.theme)

	d.helpView.SetTextAlign(tview.AlignCenter).
		SetText(`							
[#4a28ff] ██████[gray]╗[#4a28ff]  ██████[gray]╗[#4a28ff] ██████[gray]╗[#4a28ff] ██[gray]╗[#4a28ff]██[gray]╗[#4a28ff]   ██[gray]╗[#4a28ff]███████[gray]╗
[#4a28ff]██[gray]╔════╝ [#4a28ff]██[gray]╔═══[#4a28ff]██[gray]╗[#4a28ff]██[gray]╔══[#4a28ff]██[gray]╗[#4a28ff]██[gray]║[#4a28ff]██[gray]║   [#4a28ff]██[gray]║[#4a28ff]██[gray]╔════╝
[#4a28ff]██[gray]║  [#4a28ff]███[gray]╗[#4a28ff]██[gray]║   [#4a28ff]██[gray]║[#4a28ff]██████[gray]╔╝[#4a28ff]██[gray]║[#4a28ff]██[gray]║   [#4a28ff]██[gray]║[#4a28ff]███████[gray]╗
[#4a28ff]██[gray]║   [#4a28ff]██[gray]║[#4a28ff]██[gray]║   [#4a28ff]██[gray]║[#4a28ff]██[gray]╔══[#4a28ff]██[gray]╗[#4a28ff]██[gray]║[#4a28ff]██[gray]║   [#4a28ff]██[gray]║╚════[#4a28ff]██[gray]║
[gray]╚[#4a28ff]██████[gray]╔╝╚[#4a28ff]██████[gray]╔╝[#4a28ff]██████[gray]╔╝[#4a28ff]██[gray]║╚[#4a28ff]██████[gray]╔╝[#4a28ff]███████[gray]║
[gray] ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝ ╚═════╝ ╚══════╝[white]

[::]                                [yellow]v1.0.0[white]                                    

[::]                      Copyright© 2025 Gobius Developers
                                                                      
 Help Information:
 • Press '1' - '3' to switch between views
 • Press 'ctrl+c' or'q' to quit
 • Press 'v' in dashboard to change GPU view mode
 • Press complete digits of pi to get a surprise
`)

	// Initialize other components as before
	d.gpuViewer = NewGPUViewer(d.app, d.theme)

	// Initialize Log Viewer
	d.logViewer = NewLogsViewer(d.theme)

	consoleWriter := zerolog.ConsoleWriter{Out: d.logViewer, TimeFormat: "15:04:05.000000000"}
	d.logger = zerolog.New(consoleWriter).With().Timestamp().Logger()

	// Initialize Metrics View
	//d.metricsView = NewCustomTextView(d.theme)
	d.metricsView = NewMetricsView(d.theme)
	d.metricsView.
		SetTitle("Metrics").
		SetTitleAlign(tview.AlignRight).
		SetBorder(true).
		SetTitleColor(d.theme.Colors.Primary).
		SetBorderColor(d.theme.Colors.Border)

	// Initialize content area
	d.contentArea = tview.NewFlex()
	d.contentArea.SetBackgroundColor(d.theme.Colors.Background)

	d.utilizationMonitor = NewUtilizationMonitor(d.theme, 60, 1) // Height will be set dynamically
}

func (d *Dashboard) updateMenu() {
	getMenuStyle := func(isSelected bool) string {
		if isSelected {
			// Selected item: MenuFocus background with MenuText foreground
			return fmt.Sprintf("#%06x:#%06x",
				d.theme.Colors.MenuText.Hex(),
				d.theme.Colors.MenuFocus.Hex())
		}
		// Unselected item: MenuText foreground with Surface background
		return fmt.Sprintf("#%06x:#%06x",
			d.theme.Colors.MenuText.Hex(),
			d.theme.Colors.Secondary.Hex())
	}

	// Create view selection menu items with visual indicators
	leftMenuText := fmt.Sprintf(
		"[%s] 1 Dashboard [%s] [%s] 2 Logs [%s] [%s] 3 Help [%s]",
		getMenuStyle(d.currentView == ViewDashboard),
		getMenuStyle(false),
		getMenuStyle(d.currentView == ViewLogs),
		getMenuStyle(false),
		getMenuStyle(d.currentView == ViewHelp),
		getMenuStyle(false),
	)

	// Create command menu item with different styling
	rightMenuText := fmt.Sprintf("[#%06x:#%06x:b] Q Quit [-:-]",
		d.theme.Colors.Error.Hex(), // Use error color for commands
		d.theme.Colors.Secondary.Hex())

	// Set the menu texts
	d.leftMenu.SetText(leftMenuText)
	d.rightMenu.SetText(rightMenuText)
}

func (d *Dashboard) switchToView(view ViewMode) {
	d.currentView = view
	d.updateMenu() // Call updateMenu to refresh the menu display
	d.contentArea.Clear()

	switch view {
	case ViewDashboard:
		// Create the top section with GPUViewer and MetricsView
		topSection := tview.NewFlex().
			AddItem(d.gpuViewer, 0, 2, true).
			AddItem(d.metricsView, 0, 1, false)

		// Calculate minimum height to ensure we never go below 3 rows
		_, _, _, screenHeight := d.layout.GetRect()
		minHeight := max(3, screenHeight/4)
		d.utilizationMonitor.height = minHeight

		// Use flex proportions: 3 for top section, 1 for utilization (25%)
		d.contentArea.SetDirection(tview.FlexRow).
			AddItem(topSection, 0, 3, true).
			AddItem(d.utilizationMonitor, 0, 1, false)

		// Set focus to the GPUViewer
		d.app.SetFocus(d.gpuViewer)
	case ViewLogs:
		d.contentArea.AddItem(d.logViewer, 0, 1, true)
		d.app.SetFocus(d.logViewer)
	case ViewHelp:
		d.contentArea.AddItem(d.helpView, 0, 1, true)
		d.app.SetFocus(d.helpView)
	}
}

func (d *Dashboard) setupLayout() {
	d.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	d.layout.SetBackgroundColor(d.theme.Colors.Background)

	// Add title bar at the top, content in the middle, menu at the bottom
	d.layout.
		AddItem(d.titleBar, 1, 0, false).
		AddItem(d.contentArea, 0, 1, true).
		AddItem(d.menu, 1, 0, false)

	d.switchToView(ViewDashboard)

	// Warning modal for small screen size
	warningText := tview.NewTextView().
		SetText("⚠️ Terminal size too small! Please resize the window.").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	modal := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(warningText, 3, 1, false)

	_ = modal
	// Flag to prevent repeated root setting
	//isSmall := false
	// Detect screen size and show/hide warning
	d.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		//width, height := screen.Size()
		//// minWidth := 50
		// minHeight := 15

		//d.logger.Info().Msgf("🔥 Info log with color %d, %d", width, height)

		// if width < minWidth || height < minHeight {
		// 	if !isSmall {
		// 		isSmall = true
		// 		fmt.Println("Setting root to modal")
		// 		d.app.QueueUpdateDraw(func() { d.app.oot(modal, true) })
		// 	}
		// } else {
		// 	if isSmall {
		// 		isSmall = false
		// 		fmt.Println("Setting root to layout")
		// 		//d.app.QueueUpdateDraw(func() { d.app.SetRoot(d.layout, true) })
		// 	}
		// }

		return false
	})
	// Set up global keyboard handling
	// Get the existing input capture if any
	existingCapture := d.app.GetInputCapture()
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle our specific keys
		switch event.Rune() {
		case '1':
			d.switchToView(ViewDashboard)
			return nil
		case '2':
			d.switchToView(ViewLogs)
			return nil
		case '3':
			d.switchToView(ViewHelp)
			return nil
		case 'q', 'Q':
			d.app.Stop()
			return nil
		}

		// Pass to existing handler if we didn't handle it
		if existingCapture != nil {
			return existingCapture(event)
		}

		return event
	})

	d.app.SetRoot(d.layout, true)
}

// Add a function to check screen size and show warning if needed
func (d *Dashboard) checkScreenSize() {
	// Create a simple message for small screens
	smallScreenMsg := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(fmt.Sprintf(
			"\n\n[red]Screen too small![white]\n\nPlease resize your terminal to at least %dx%d\n\nPress 'q' to quit",
			MinScreenWidth, MinScreenHeight))

	// Add a border
	msgBox := tview.NewFrame(smallScreenMsg).
		SetBorders(1, 1, 1, 1, 1, 1).
		SetBorderColor(tcell.ColorRed)

	// Create a simple modal page that shows when screen is too small
	smallScreenPage := tview.NewPages().
		AddPage("small", msgBox, true, true)

	// Add input handling to the small screen page
	smallScreenPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			d.app.Stop()
		}
		return event
	})

	// Create a pages primitive to switch between normal and small screen
	pages := tview.NewPages().
		AddPage("normal", d.layout, true, true).
		AddPage("small", smallScreenPage, true, false)

	// Set the root to our pages container
	d.app.SetRoot(pages, true)

	// Add a simple resize handler
	d.app.SetFocus(d.layout)

	// Create a primitive that will be used to detect screen size
	sizeDetector := tview.NewBox().
		SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			// Check if screen is too small
			if width < MinScreenWidth || height < MinScreenHeight {
				// Switch to small screen page
				pages.SwitchToPage("small")
			} else {
				// Switch to normal page
				pages.SwitchToPage("normal")
			}
			return x, y, width, height
		})

	// Add the size detector to the layout directly
	d.layout.AddItem(sizeDetector, 0, 0, false)
}

func (d *Dashboard) startUpdateHandler() {
	go func() {
		for update := range d.updates {
			d.app.QueueUpdateDraw(func() {
				switch update.Type {
				case UpdateGPUs:
					if gpus, ok := update.Payload.([]metrics.GPUInfo); ok {
						d.gpuViewer.UpdateGPUs(gpus)
						d.utilizationMonitor.AddSample(gpus)
					}
				case UpdateLog:
					if log, ok := update.Payload.(string); ok {
						fmt.Fprintf(d.logViewer, "%s", log)
					}
				case UpdateValidatorMetrics:
					if metrics, ok := update.Payload.(ValidatorMetrics); ok {
						d.mu.Lock()
						d.validatorMetrics = metrics
						d.metricsView.UpdateMetrics(d.validatorMetrics, d.financialMetrics)
						d.mu.Unlock()
					}
				case UpdateFinancialMetrics:
					if metrics, ok := update.Payload.(FinancialMetrics); ok {
						d.mu.Lock()
						d.financialMetrics = metrics
						d.metricsView.UpdateMetrics(d.validatorMetrics, d.financialMetrics)
						d.mu.Unlock()
					}
				}
			})
		}
	}()
}

// Add this method to the MetricsView struct
func (m *MetricsView) UpdateMetrics(validatorMetrics ValidatorMetrics, financialMetrics FinancialMetrics) {
	var content strings.Builder

	// Validator Metrics
	content.WriteString(fmt.Sprintf("[#%06x::b]Validator Metrics[-:-:-]\n", m.theme.Colors.Primary.Hex()))
	content.WriteString(fmt.Sprintf("%s Session Time: %s\n", string(m.theme.Symbols.Bullet), validatorMetrics.SessionTime))
	content.WriteString(fmt.Sprintf("%s Solutions: [#%06x]%d/%d (%.1f%%)[-:-:-]\n",
		string(m.theme.Symbols.Bullet),
		m.theme.Colors.Success.Hex(),
		validatorMetrics.SolutionsLastMinute.Success,
		validatorMetrics.SolutionsLastMinute.Total,
		validatorMetrics.SolutionsLastMinute.Rate*100))
	content.WriteString(fmt.Sprintf("%s Avg Solutions/min: %.2f\n\n", string(m.theme.Symbols.Bullet), validatorMetrics.AverageSolutionsPerMin))

	// Financial Metrics
	content.WriteString(fmt.Sprintf("[#%06x::b]Financial Metrics[-:-:-]\n", m.theme.Colors.Secondary.Hex()))
	content.WriteString(fmt.Sprintf("%s Token/min: %.4f\n", string(m.theme.Symbols.Bullet), financialMetrics.TokenIncomePerMin))
	content.WriteString(fmt.Sprintf("%s USD/hour: $%.2f\n", string(m.theme.Symbols.Bullet), financialMetrics.IncomePerHour))
	content.WriteString(fmt.Sprintf("%s Profit/day: [#%06x]$%.2f[-:-:-]",
		string(m.theme.Symbols.Bullet),
		m.theme.Colors.Success.Hex(),
		financialMetrics.ProfitPerDay))

	m.SetText(content.String())
}

// Demo main function
func main() {
	dashboard := NewDashboard()
	//dashboard.checkScreenSize() // Add screen size check

	//gpuStates := []string{"Mining", "Error", "Idle"}

	// Redirect Zerolog output to the log view

	var dummyGpus []*metrics.GPUInfo

	solveTimes := []float64{5.2, 4.8, 5.6, 5.0, 7.0, 4.0}

	for i := 0; i < 64; i++ {
		status := "Idle"
		// convert float64 to seconds and round to 1 decimal place
		solveTime := solveTimes[i%len(solveTimes)]
		duration := time.Duration(solveTime * float64(time.Second))
		dummyGpus = append(dummyGpus, &metrics.GPUInfo{
			ID:         i,
			Status:     status,
			SolveTime:  duration,
			MaxSolves:  100,
			SolveTimes: make([]float64, 0, 100),
		})
	}
	// launch a goroutine per gpu to fake mining and update the gpu status
	for id, gpu := range dummyGpus {
		go func(gpu *metrics.GPUInfo, gpuId int) {
			i := gpuId
			errorCount := 0
			for {
				gpu.UpdateStatus("Mining")
				rndSolveTime := solveTimes[i%len(solveTimes)] + rand.Float64()*0.5 - 0.25
				time.Sleep(time.Duration(rndSolveTime * float64(time.Second)))
				// everyone so often randomly "error"
				if rand.Float64() < 0.1 {
					gpu.UpdateStatus("Error")
					time.Sleep(time.Duration(10 * time.Second))
					errorCount++
				} else {
					gpu.UpdateSolveTime(time.Duration(rndSolveTime * float64(time.Second)))
					gpu.AddSolveTime(rndSolveTime)
					gpu.UpdateStatus("Idle")
					// every so often we will stay idel for some seconds:
					if rand.Float64() < 0.1 {
						time.Sleep(time.Duration(rand.Float64()*20+3) * time.Second)
					}
					//errorCount = 0
				}
				i++
				if errorCount > 3 {
					gpu.UpdateStatus("Error")
					return
				}
			}
		}(gpu, id)
	}
	// Simulate updates
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		logger := dashboard.GetLogger()

		count := 0

		for {
			select {
			case <-ticker.C:

				// Update GPUs
				var gpus []metrics.GPUInfo
				for i := 0; i < 64; i++ {
					gpu := dummyGpus[i]
					// // add a random +/- solve time to the gpu
					// rndSolveTime := solveTimes[i%len(solveTimes)] + rand.Float64()*0.5 - 0.25
					// gpu.AddSolveTime(rndSolveTime)
					//status := gpuStates[(count+i)%len(gpuStates)]
					gpu.ReadLock()
					gpus = append(gpus, metrics.GPUInfo{
						ID:         gpu.ID,
						Status:     gpu.Status,
						SolveTime:  gpu.SolveTime,
						SolveTimes: gpu.SolveTimes,
						MaxSolves:  gpu.MaxSolves,
					})
					gpu.ReadUnlock()
				}
				count++

				dashboard.updates <- StateUpdate{
					Type:    UpdateGPUs,
					Payload: gpus,
				}

				// // Update logs
				// dashboard.updates <- StateUpdate{
				// 	Type: UpdateLog,
				// 	Payload: fmt.Sprintf("[#%06x]%s[white] New log entry at %s\n",
				// 		dashboard.theme.Colors.Primary.Hex(),
				// 		string(dashboard.theme.Symbols.Bullet),
				// 		time.Now().Format("15:04:05")),
				// }
				logger.Info().Msgf("New log entry at %s", time.Now().Format("15:04:05"))

				// Update metrics
				dashboard.updates <- StateUpdate{
					Type: UpdateValidatorMetrics,
					Payload: ValidatorMetrics{
						SessionTime:      time.Hour.String(),
						SolvedLastMinute: int64(count % 10),
					},
				}

				// Update financial metrics
				dashboard.updates <- StateUpdate{
					Type: UpdateFinancialMetrics,
					Payload: FinancialMetrics{
						TokenIncomePerMin:  float64(count) * 0.1,
						TokenIncomePerHour: float64(count) * 6.0,
						TokenIncomePerDay:  float64(count) * 144.0,
						IncomePerMin:       float64(count) * 0.5,
						IncomePerHour:      float64(count) * 30.0,
						IncomePerDay:       float64(count) * 720.0,
						ProfitPerMin:       float64(count) * 0.3,
						ProfitPerHour:      float64(count) * 18.0,
						ProfitPerDay:       float64(count) * 432.0,
					},
				}
			}
		}
	}()

	if err := dashboard.app.Run(); err != nil {
		panic(err)
	}
}

func generateSparkline(solveTimes []float64, width int) string {
	if len(solveTimes) == 0 {
		return ""
	}

	// Calculate average
	var sum float64
	for _, t := range solveTimes {
		sum += t
	}
	avg := sum / float64(len(solveTimes))

	// Build colored sparkline
	var sparkline strings.Builder

	// Define our sparkline character
	const spark = "▄"

	// Color ranges (as percentages of average)
	colors := []struct {
		threshold float64
		color     tcell.Color
	}{
		{1.5, tcell.NewRGBColor(255, 0, 0)},    // Red (50% above average)
		{1.25, tcell.NewRGBColor(255, 128, 0)}, // Orange (25% above average)
		{1.1, tcell.NewRGBColor(255, 192, 0)},  // Yellow-Orange (10% above average)
		{1.0, tcell.NewRGBColor(255, 255, 0)},  // Yellow (at average)
		{0.9, tcell.NewRGBColor(192, 255, 0)},  // Yellow-Green (10% below average)
		{0.75, tcell.NewRGBColor(128, 255, 0)}, // Light Green (25% below average)
		{0.0, tcell.NewRGBColor(0, 255, 0)},    // Green (faster)
	}

	for i := max(0, len(solveTimes)-width); i < len(solveTimes); i++ {
		ratio := solveTimes[i] / avg

		// Find appropriate color
		var color tcell.Color
		for _, c := range colors {
			if ratio >= c.threshold {
				color = c.color
				break
			}
		}

		// Add colored sparkline character
		sparkline.WriteString(fmt.Sprintf("[#%06x]%s", color.Hex(), spark))
	}

	sparkline.WriteString("[-]") // Reset color
	return sparkline.String()
}

func (g *GPUViewer) UpdateGPUs(gpus []metrics.GPUInfo) {
	// Store current selection and scroll state before updating
	// prevRow, prevCol := g.table.GetSelection()
	// // Save the current scroll position
	// offsetRow, offsetCol := g.table.GetOffset()
	g.gpus = gpus

	// Only apply threshold check if we're in DetailedView and this is the first update
	if len(g.gpus) == 0 && len(gpus) > g.compactThreshold {
		g.viewMode = CompactView
	}

	switch g.viewMode {
	case DetailedView:
		g.renderDetailedView()
	case CompactView:
		g.renderCompactView()
	case SummaryView:
		g.renderSummaryView()
	}

	// Reapply the previous selection if it still exists. For DetailedView, ensure the header row is skipped.
	// rowToSelect := prevRow
	// if g.viewMode == DetailedView && rowToSelect < 1 {
	// 	rowToSelect = 1
	// }
	// if g.table.GetCell(rowToSelect, 0) == nil {
	// 	if g.viewMode == DetailedView {
	// 		rowToSelect = 1
	// 	} else {
	// 		rowToSelect = 0
	// 	}
	// }

	//g.table.SetOffset(offsetRow, offsetCol)
	// Set the selection, which may trigger auto-scrolling
	//g.table.Select(rowToSelect, prevCol)
	// Manually restore scroll position

}
