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
	viewPtr        atomic.Pointer[CustomTextView] // Use atomic pointer
	writerPtr      atomic.Pointer[io.Writer]      // Use atomic pointer
	logCh          chan []byte
	originalStdout io.Writer
	originalStderr io.Writer
}

func NewLogRouter() (*logrouter, func()) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	router := &logrouter{
		// viewPtr and writerPtr start as nil implicitly
		originalStdout: originalStdout,
		originalStderr: originalStderr,
		logCh:          make(chan []byte, 16*4096),
	}
	router.Start() // Start the log channel processor early

	// Create pipe and replace stdout/stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// Start Pipe Reader Goroutine
	exit := make(chan bool)
	go func() {
		defer close(exit)
		buffer := make([]byte, 4096) // Read buffer
		for {
			n, err := r.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(originalStderr, "LogRouter pipe read error: %v\n", err)
				}
				break // Exit loop on EOF or error
			}
			if n > 0 {
				dataToWrite := buffer[:n]
				// Check view and headless state atomically
				view := router.viewPtr.Load()

				// If view is nil (never set or stopped) or headless, write directly
				if view == nil {
					_, _ = router.originalStderr.Write(dataToWrite)
				} else {
					// Otherwise, route via the Write method for TUI channel queuing
					_, _ = router.Write(dataToWrite)
				}
			}
		}
	}()

	// Cleanup function
	cleanup := func() {
		w.Close() // Close pipe writer, signals EOF to reader goroutine
		<-exit    // Wait for reader goroutine to finish processing pipe
		close(router.logCh)
		// Restore original stdout/stderr
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}

	return router, cleanup
}

// SetView is called once when the TUI is ready
func (tw *logrouter) SetView(view *CustomTextView) {
	ansiWriter := tview.ANSIWriter(view)
	tw.viewPtr.Store(view)
	tw.writerPtr.Store(&ansiWriter) // Store address of the writer
}

// StopTUI reverts logging back to the original console stderr
func (tw *logrouter) Stop() {
	tw.viewPtr.Store(nil)
	tw.writerPtr.Store(nil)
}

// Write is now only called when routing to TUI (view != nil and not headless)
func (tw *logrouter) Write(p []byte) (n int, err error) {
	logCopy := make([]byte, len(p))
	copy(logCopy, p)
	select {
	case tw.logCh <- logCopy:
		return len(p), nil
	default:
		warning := fmt.Sprintf("[WARN] LogRouter channel full, TUI message dropped: %s\n", string(p))
		tw.originalStderr.Write([]byte(warning))
		return tw.originalStderr.Write(p)
	}
}

// Start runs the goroutine that processes the log channel for the TUI view
func (tw *logrouter) Start() {
	go func() {
		for msg := range tw.logCh {
			// Load view and writer atomically
			view := tw.viewPtr.Load()
			writerAddr := tw.writerPtr.Load()

			// If view is set (StopTUI not called) and not initially headless
			if view != nil && writerAddr != nil {
				// Dereference the writer pointer
				writer := *writerAddr
				fmt.Fprintf(writer, "%s", msg)
				// Still need QueueUpdateDraw here if relying on this path
				// (Let's omit it for now as per user preference, but be aware)
			} else {
				// Fallback if view is nil (stopped)
				tw.originalStderr.Write(msg)
			}
		}
	}()
}
