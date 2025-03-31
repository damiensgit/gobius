package tui

import (
	"github.com/rivo/tview"
)

// LogViewer component
type LogViewer struct {
	*CustomTextView
	theme    *Theme
	maxLines int
}

// func (w *LogViewer) Write(p []byte) (n int, err error) {

// 	if w.view == nil {
// 		return os.Stderr.Write(p)
// 	}

// 	// Ensure logs appear immediately in the TextView
// 	fmt.Fprintf(w.view.TextView, "%s", p)

// 	return len(p), nil
// 	// return tview.ANSIWriter(w.TextView).Write(p)
// }

func NewLogsViewer(theme *Theme) *LogViewer {

	logsView := NewCustomTextView(theme)
	logsView.
		SetTitle(" Logs ").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitleColor(theme.Colors.Primary).
		SetBorderColor(theme.Colors.Border).
		SetScrollable(true).
		ScrollToEnd()

	return &LogViewer{
		CustomTextView: logsView,
		theme:          theme,
		maxLines:       1000,
	}
}
