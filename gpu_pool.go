package main

import (
	"context"
	"fmt"
	task "gobius/common"
	"gobius/config"
	"gobius/metrics"
	"gobius/models"
	"gobius/utils"
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// GPUPool manages a collection of GPU clients.
type GPUPool struct {
	mu      sync.RWMutex
	gpus    []*task.GPU
	config  *config.AppConfig
	logger  zerolog.Logger
	gpura   *utils.RunningAverage // Running average for solve times across all GPUs
	modelID string
}

// NewGPUPool creates and initializes a new GPUPool.
func NewGPUPool(cfg *config.AppConfig, logger zerolog.Logger, modelID string, mockGPUs int) (*GPUPool, error) {
	ra := utils.NewRunningAverage(15 * time.Minute) // TODO: Make configurable?
	pool := &GPUPool{
		config:  cfg,
		logger:  logger,
		gpura:   ra,
		modelID: modelID,
	}

	gpusURLS, ok := cfg.ML.Cog[modelID]
	if !ok {
		return nil, fmt.Errorf("missing GPU URLs for model %s in config", modelID)
	}

	gpuList := zerolog.Arr()
	for id, gpuUrl := range gpusURLS.URL {
		gpuList.Str(gpuUrl)
		pool.gpus = append(pool.gpus, task.NewGPU(id, gpuUrl))
	}
	logger.Info().Array("gpus", gpuList).Msg("initialized GPUs from config")

	// Add mock GPUs if requested
	id := len(pool.gpus)
	for i := 0; i < mockGPUs; i++ {
		gpu := task.NewGPU(id, "")
		gpu.Mock = true
		pool.gpus = append(pool.gpus, gpu)
		id++
	}
	if mockGPUs > 0 {
		logger.Warn().Int("count", mockGPUs).Msg("added mock GPUs")
	}

	if len(pool.gpus) == 0 {
		logger.Warn().Msg("no GPUs configured or mocked")
	}

	return pool, nil
}

// ValidateGPUs checks the health and compatibility of the configured GPUs with the model.
// It returns an error if all GPUs fail validation, or if no GPUs are configured.
// If mock GPUs are used, it will return a warning and not an error.
// It does not lock the GPU pool, so should be called early in the startup process before the pool is used.
func (p *GPUPool) ValidateGPUs(model models.ModelInterface) error {
	if len(p.gpus) == 0 {
		p.logger.Warn().Msg("no GPUs to validate")
		return nil
	}

	p.logger.Info().Str("model", model.GetID()).Msg("validating model on gpu(s)")

	mu := sync.Mutex{}
	mockGPUs := 0
	fastestTimeSeen := time.Duration(0)
	times := make([]time.Duration, len(p.gpus))
	var wg sync.WaitGroup

	for i, gpu := range p.gpus {
		if gpu.Mock {
			mockGPUs++
			continue
		}
		wg.Add(1)
		go func(i int, gpu *task.GPU) {
			defer wg.Done()
			start := time.Now()
			// Use a unique task ID for validation to avoid cache issues if any
			validationTaskID := fmt.Sprintf("startup-test-taskid-gpu-%d-%d", gpu.ID, rand.Intn(10000))
			err := model.Validate(gpu, validationTaskID)
			timeTaken := time.Since(start)

			if err != nil {
				// Use Fatalf for critical startup errors
				p.logger.Error().Err(err).Int("gpu", gpu.ID).Str("url", gpu.Url).Msg("error validating the model on GPU")
				gpu.Enabled = false
			} else {
				p.logger.Info().Int("gpu", gpu.ID).Str("duration", timeTaken.String()).Msg("GPU validation successful")
				times[i] = timeTaken
				mu.Lock()
				if fastestTimeSeen == 0 || timeTaken < fastestTimeSeen {
					fastestTimeSeen = timeTaken
				}
				mu.Unlock()
			}
		}(i, gpu)
	}

	wg.Wait()

	// Filter out times from GPUs that failed validation (time will be 0)
	validTimes := []time.Duration{}
	for _, t := range times {
		if t > 0 {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) > 0 {
		averageTime := averageDurations(validTimes)
		p.logger.Info().Str("average", averageTime.String()).Str("fastest", fastestTimeSeen.String()).Int("validated", len(validTimes)).Int("total", len(p.gpus)).Msg("GPU validation complete")
	} else if len(p.gpus) > 0 {
		if mockGPUs == 0 {
			p.logger.Error().Msg("all GPUs failed validation")
			return fmt.Errorf("all configured GPUs failed validation for model %s", p.modelID)
		} else {
			p.logger.Warn().Msg("all GPUs failed validation, but mock GPUs were used")
		}
	} else {
		// This case was handled at the start, but included for completeness
		p.logger.Warn().Msg("no GPUs were configured for validation")
	}

	return nil
}

// GetGPUs returns a copy of the list of GPU clients.
func (p *GPUPool) GetGPUs() []*task.GPU {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// Return a copy to prevent external modification of the internal slice
	gpusCopy := make([]*task.GPU, len(p.gpus))
	copy(gpusCopy, p.gpus)
	return gpusCopy
}

// NumGPUs returns the total number of GPUs managed by the pool.
func (p *GPUPool) NumGPUs() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.gpus)
}

// NumEnabledGPUs returns the number of GPUs currently enabled.
func (p *GPUPool) NumEnabledGPUs() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	count := 0
	for _, gpu := range p.gpus {
		if gpu.IsEnabled() {
			count++
		}
	}
	return count
}

// AddSolveTime adds a solve duration to the running average.
func (p *GPUPool) AddSolveTime(d time.Duration) {
	p.gpura.Add(d)
}

// AverageSolveTime returns the average solve time.
func (p *GPUPool) AverageSolveTime() time.Duration {
	return p.gpura.Average()
}

// GetGPUInfoForMetrics returns GPU information suitable for metrics/TUI display.
func (p *GPUPool) GetGPUInfoForMetrics() []metrics.GPUInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	avgSolveTime := p.AverageSolveTime()
	gpuInfos := make([]metrics.GPUInfo, 0, len(p.gpus))
	for _, gpu := range p.gpus {
		status := gpu.Status // Use the thread-safe getter
		info := metrics.GPUInfo{
			ID:        gpu.ID,
			Status:    status,
			SolveTime: avgSolveTime, // Show pool average for now
			// TODO: Add SolveTimes and MaxSolves if needed
		}
		if !gpu.IsEnabled() {
			info.Status = "Error" // Override status if disabled
		}
		gpuInfos = append(gpuInfos, info)
	}
	return gpuInfos
}

// averageDurations calculates the average of a slice of time.Durations.
// Handles empty or nil slices gracefully.
func averageDurations(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}
	total := time.Duration(0)
	for _, t := range times {
		total += t
	}
	return total / time.Duration(len(times))
}

// StartStatusResetter periodically resets the error state of GPUs.
func (p *GPUPool) StartStatusResetter(ctx context.Context, resetInterval time.Duration) {
	ticker := time.NewTicker(resetInterval)
	defer ticker.Stop()

	p.logger.Info().Dur("interval", resetInterval).Msg("starting GPU status resetter")

	for {
		select {
		case <-ctx.Done():
			p.logger.Info().Msg("stopping GPU status resetter")
			return
		case <-ticker.C:
			p.mu.Lock() // Full lock needed to modify GPU state
			p.logger.Debug().Msg("attempting to reset GPU error states")
			for _, gpu := range p.gpus {
				if !gpu.IsEnabled() {
					p.logger.Info().Int("gpu", gpu.ID).Msg("resetting error state for GPU")
					gpu.ResetErrorState()
				}
			}
			p.mu.Unlock()
		}
	}
}
