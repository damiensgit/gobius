package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
