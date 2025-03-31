package tui

import (
	"github.com/rivo/tview"
)

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
