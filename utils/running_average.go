package utils

import (
	"sync"
	"time"
)

type DurationRecord struct {
	Time     time.Time
	Duration time.Duration
}

type RunningAverage struct {
	mu        sync.RWMutex
	Durations []DurationRecord
	Sum       time.Duration
	Interval  time.Duration
}

func NewRunningAverage(interval time.Duration) *RunningAverage {
	return &RunningAverage{
		Durations: make([]DurationRecord, 0),
		Sum:       0,
		Interval:  interval,
	}
}

func (ra *RunningAverage) Add(duration time.Duration) {
	now := time.Now()

	ra.mu.Lock()
	defer ra.mu.Unlock()

	// Add the new duration
	ra.Durations = append(ra.Durations, DurationRecord{Time: now, Duration: duration})
	ra.Sum += duration

	// Remove durations older than the interval
	for len(ra.Durations) > 0 && now.Sub(ra.Durations[0].Time) > ra.Interval {
		ra.Sum -= ra.Durations[0].Duration
		ra.Durations = ra.Durations[1:]
	}
}

func (ra *RunningAverage) Average() time.Duration {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	if len(ra.Durations) == 0 {
		return 0
	}
	return ra.Sum / time.Duration(len(ra.Durations))
}
