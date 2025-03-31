package tui

import (
	"time"

	"github.com/gdamore/tcell/v2"
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
