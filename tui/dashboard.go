package tui

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
	LogViewer          *LogViewer
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

	return d
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

[::][yellow]v1.0.0[white]                                    

[::]Copyright© 2025 Gobius Developers
[::]Portions Copyright© 2025 Arbius Developers

                                                                      
 Help Information:
 • Press '1' - '3' to switch between views
 • Press 'ctrl+c' or'q' to quit
 • Press 'v' in dashboard to change GPU view mode
 • Press complete digits of pi to get a surprise
`)

	// Initialize other components as before
	d.gpuViewer = NewGPUViewer(d.app, d.theme)

	// Initialize Log Viewer
	d.LogViewer = NewLogsViewer(d.theme)

	// consoleWriter := zerolog.ConsoleWriter{Out: d.LogViewer, TimeFormat: "15:04:05.000000000"}
	// d.logger = zerolog.New(consoleWriter).With().Timestamp().Logger()

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
		d.contentArea.AddItem(d.LogViewer, 0, 1, true)
		d.app.SetFocus(d.LogViewer)
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
		width, height := screen.Size()
		// minWidth := 50
		// minHeight := 15

		d.logger.Info().Msgf("🔥 Info log with color %d, %d", width, height)

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
				// case UpdateLog:
				// 	if log, ok := update.Payload.(string); ok {
				// 		fmt.Fprintf(d.LogViewer, "%s", log)
				// 	}
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

func (d *Dashboard) Run() {
	d.startUpdateHandler()

	if err := d.app.Run(); err != nil {
		panic(err)
	}
}

func (d *Dashboard) Stop() {

	if d.app == nil {
		panic("dashboard stop before initialize")
	}

	d.app.Stop()
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

				// Update logs
				dashboard.updates <- StateUpdate{
					Type: UpdateLog,
					Payload: fmt.Sprintf("[#%06x]%s[white] New log entry at %s\n",
						dashboard.theme.Colors.Primary.Hex(),
						string(dashboard.theme.Symbols.Bullet),
						time.Now().Format("15:04:05")),
				}

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
