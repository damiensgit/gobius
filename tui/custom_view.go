package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
