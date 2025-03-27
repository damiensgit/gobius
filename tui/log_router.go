package tui

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/rivo/tview"
)

// basic log output router to direct logging to either console or tui if avail
type logrouter struct {
	view           *CustomTextView
	Headless       atomic.Bool
	writer         io.Writer
	logCh          chan []byte
	originalStdout io.Writer
	originalStderr io.Writer
}

func NewLogRouter() (*logrouter, func()) {
	// Save original values
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	router := &logrouter{
		view:           nil,
		Headless:       atomic.Bool{},
		originalStdout: originalStdout,
		originalStderr: originalStderr,
		logCh:          make(chan []byte, 16*4096),
	}
	router.Start()

	// Create a single pipe
	r, w, _ := os.Pipe()

	// Replace both stdout and stderr with the same pipe writer
	os.Stdout = w
	os.Stderr = w

	// Create exit channel
	exit := make(chan bool)

	// Start a single goroutine to handle both outputs
	go func() {
		io.Copy(router, r)
		exit <- true
	}()

	// Return router and cleanup function
	return router, func() {
		w.Close()
		<-exit
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}
}

func (tw *logrouter) SetView(view *CustomTextView) {
	tw.view = view
	tw.writer = tview.ANSIWriter(view)
}

func (tw *logrouter) Write(p []byte) (n int, err error) {
	isHeadless := tw.Headless.Load()
	if tw.view == nil || isHeadless {
		return tw.originalStderr.Write(p)
	} else {
		// Queue log message in channel
		logCopy := make([]byte, len(p))
		copy(logCopy, p)
		select {
		case tw.logCh <- logCopy:
			return len(p), nil
		default:
			// Channel full, fallback to stderr
			return tw.originalStderr.Write(p)
		}
	}
}

func (tw *logrouter) Start() {
	go func() {
		for msg := range tw.logCh {
			// Re-check state before writing to view
			if tw.view != nil && !tw.Headless.Load() {
				fmt.Fprintf(tw.writer, "%s", msg)
			} else {
				// State changed while message was in queue
				tw.originalStderr.Write(msg)
			}
		}
	}()
}
