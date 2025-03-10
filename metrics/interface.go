package metrics

import (
	"sync"
	"time"
)

type MetricsProvider interface {
	// Session metrics
	GetSessionTime() string
	GetSolvedLastMinute() int64
	GetSuccessCount() int64
	GetTotalCount() int64
	GetSuccessRate() float64
	GetAverageSolutionRate() float64
	GetAverageSolutionsPerMin() float64
	GetAverageSolvesPerMin() float64

	// Financial metrics
	GetTokenIncomePerMin() float64
	GetTokenIncomePerHour() float64
	GetTokenIncomePerDay() float64
	GetIncomePerMin() float64
	GetIncomePerHour() float64
	GetIncomePerDay() float64
	GetProfitPerMin() float64
	GetProfitPerHour() float64
	GetProfitPerDay() float64
}

type HardwareMetrics interface {
	// GPU metrics
	GetGPUs() []GPUInfo
}

type GPUInfo struct {
	ID         int
	Status     string
	SolveTime  time.Duration
	SolveTimes []float64 // Store last n solve times
	MaxSolves  int       // Maximum number of solve times to store
	mu         sync.RWMutex
}

// update status:
func (g *GPUInfo) UpdateStatus(status string) {
	g.mu.Lock()
	g.Status = status
	g.mu.Unlock()
}

func (g *GPUInfo) UpdateSolveTime(solveTime time.Duration) {
	g.mu.Lock()
	g.SolveTime = solveTime
	g.mu.Unlock()
}

// Add method to update solve times
func (g *GPUInfo) AddSolveTime(solveTime float64) {
	g.mu.Lock()
	if len(g.SolveTimes) >= g.MaxSolves {
		// Remove oldest entry
		g.SolveTimes = g.SolveTimes[1:]
	}
	g.SolveTimes = append(g.SolveTimes, solveTime)
	g.mu.Unlock()
}

func (g *GPUInfo) ReadLock() {
	g.mu.RLock()
}

func (g *GPUInfo) ReadUnlock() {
	g.mu.RUnlock()
}
