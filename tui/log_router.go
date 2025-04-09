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
	originalStdout *os.File // Keep as *os.File
	originalStderr *os.File // Keep as *os.File
	pipeWriter     *os.File // Added to store the pipe writer for restoration checks
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
	router.pipeWriter = w // Store the writer
	os.Stdout = w
	os.Stderr = w

	// Start Pipe Reader Goroutine
	exit := make(chan bool) // Keep the exit channel for cleanup synchronization
	go func() {
		defer close(exit)
		buffer := make([]byte, 4096) // Read buffer
		for {
			n, err := r.Read(buffer)
			if err != nil {
				// Don't log error if it's EOF potentially caused by cleanup closing the writer
				if err != io.EOF {
					// Use originalStderr directly for error logging to avoid loops if pipe is broken
					fmt.Fprintf(originalStderr, "LogRouter pipe read error: %v\n", err)
				}
				break // Exit loop on EOF or error
			}
			if n > 0 {
				dataToWrite := make([]byte, n) // Make a copy for concurrent safety
				copy(dataToWrite, buffer[:n])

				view := router.viewPtr.Load()

				// If view is nil (never set or stopped), write directly to original stderr
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
		// Restore outputs first, using the idempotent method
		router.RestoreOutputs()

		// Close pipe writer, signals EOF to reader goroutine
		if router.pipeWriter != nil { // Check if already closed/nil
			router.pipeWriter.Close()
			router.pipeWriter = nil // Prevent future restoration attempts after cleanup
		}

		// Wait for reader goroutine to finish processing pipe
		<-exit

		// Close the log channel
		close(router.logCh)

		// Explicit restore removed, handled by RestoreOutputs now
	}

	return router, cleanup
}

// SetView is called once when the TUI is ready
func (tw *logrouter) SetView(view *CustomTextView) {
	ansiWriter := tview.ANSIWriter(view)
	tw.viewPtr.Store(view)
	tw.writerPtr.Store(&ansiWriter) // Store address of the writer
}

// StopTUIOutput reverts logging back to the original console stderr if TUI was active.
// It does NOT affect the pipe redirection itself.
func (tw *logrouter) StopTUIOutput() {
	tw.viewPtr.Store(nil)
	tw.writerPtr.Store(nil)
}

// RestoreOutputs safely restores original os.Stdout and os.Stderr if they are currently redirected.
// This method is idempotent and safe to call multiple times (e.g., manually and via defer cleanup).
func (tw *logrouter) RestoreOutputs() {
	// Check if redirection is active before restoring
	if tw.pipeWriter != nil && os.Stdout == tw.pipeWriter {
		os.Stdout = tw.originalStdout
		// Attempt to flush stdout after restoring
		_ = os.Stdout.Sync()
	}
	if tw.pipeWriter != nil && os.Stderr == tw.pipeWriter {
		os.Stderr = tw.originalStderr
		// Attempt to flush stderr after restoring
		_ = os.Stderr.Sync()
	}
	// We don't nil pipeWriter here, cleanup handles final closing/nilling
}

// Write directs the byte slice to the TUI log channel if the TUI is active,
// otherwise writes to the original stderr.
func (tw *logrouter) Write(p []byte) (n int, err error) {
	logCopy := make([]byte, len(p))
	copy(logCopy, p)
	select {
	case tw.logCh <- logCopy:
		return len(p), nil
	default:
		// Fallback to originalStderr if channel is full
		warning := fmt.Sprintf("[WARN] LogRouter channel full, TUI message dropped: %s\n", string(p))
		_, _ = tw.originalStderr.Write([]byte(warning)) // Use original directly
		return tw.originalStderr.Write(p)               // Use original directly
	}
}

// Start runs the goroutine that processes the log channel for the TUI view
func (tw *logrouter) Start() {
	go func() {
		for msg := range tw.logCh {
			// Load view and writer atomically
			view := tw.viewPtr.Load()
			writerAddr := tw.writerPtr.Load()

			// If view is set (StopTUIOutput not called)
			if view != nil && writerAddr != nil {
				// Dereference the writer pointer
				writer := *writerAddr
				// We still need the application draw for TUI updates
				// Assuming the caller handles QueueUpdateDraw
				fmt.Fprintf(writer, "%s", msg)
			} else {
				// Fallback if view is nil (stopped)
				// Use originalStderr directly
				_, _ = tw.originalStderr.Write(msg)
			}
		}
	}()
}
