package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
