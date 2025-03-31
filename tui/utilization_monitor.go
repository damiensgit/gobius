package tui

import (
	"gobius/metrics"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
