package tui

import (
	"fmt"
	"gobius/metrics"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
