package common

import (
	crypto "crypto/rand"
	"math/rand"
	"sync"
	"time"
)

type GPU struct {
	ID             int
	Url            string
	Enabled        bool
	ErrorCount     int
	TimeSinceError time.Time
	Mock           bool
	mu             sync.RWMutex
}

// NewGPU creates a new GPU with the given ID and URL
func NewGPU(id int, url string) *GPU {
	return &GPU{
		ID:         id,
		Url:        url,
		Enabled:    true,
		ErrorCount: 0,

		// Set the time since error to some default value so that we don't immediately disable the GPU
		TimeSinceError: time.Now(),
	}
}

const (
	maxErrorCount = 3
	errorTimeout  = 5 * time.Minute
)

// atomically increment the error count and set the time since error to now
// ensuring thread safety
func (g *GPU) IncrementErrorCount() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if time.Since(g.TimeSinceError) < errorTimeout {
		g.ErrorCount++
		if g.ErrorCount > maxErrorCount {
			g.Enabled = false
		}
	} else {
		// If the time since the last error is more than 5 minutes, reset the error count
		g.ErrorCount = 1
	}

	g.TimeSinceError = time.Now()
}

// Set enabled to true and error count to 0 ensuring thread safety once the timesinceeror is greater than 5 minutes
func (g *GPU) ResetErrorState() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Enabled = true
	g.ErrorCount = 0
	g.TimeSinceError = time.Now()

}

// get enabled flag thread safety
func (g *GPU) IsEnabled() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.Enabled
}

func (g *GPU) GetMockCid(taskid string, input interface{}) ([]byte, error) {

	b := make([]byte, 34)
	_, err := crypto.Read(b)
	if err != nil {
		return nil, err
	}
	// random sleep between 4.8 seconds and 5.5
	time.Sleep(time.Duration(rand.Intn(700)+4800) * time.Millisecond)

	return b, nil
}
