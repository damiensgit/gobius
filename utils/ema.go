package utils

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type PriceRecord struct {
	Time  time.Time
	Price float64
}

type MovingAveragePrice struct {
	mu       sync.RWMutex
	Prices   []PriceRecord
	Max      float64
	Min      float64
	EMA      float64
	Alpha    float64
	Interval time.Duration
}

// NewMovingAveragePrice creates a new MovingAveragePrice with the given number of periods and interval
// N is the number of periods to consider for the EMA calculation. E.g. N=10 will consider the last 10 periods.
func NewMovingAveragePrice(N int, interval time.Duration) *MovingAveragePrice {
	alpha := 2 / float64(N+1)
	return &MovingAveragePrice{
		Prices:   make([]PriceRecord, 0),
		Max:      math.Inf(-1),
		Min:      math.Inf(1),
		EMA:      0,
		Alpha:    alpha,
		Interval: interval,
	}
}

func (ma *MovingAveragePrice) String() string {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	if len(ma.Prices) == 0 {
		return "warming up"
	}
	return fmt.Sprintf("max: %f, min: %f, ema: %f", ma.Max, ma.Min, ma.EMA)
}

func (ma *MovingAveragePrice) Add(price float64) {
	now := time.Now()

	ma.mu.Lock()
	defer ma.mu.Unlock()

	// Add the new price
	ma.Prices = append(ma.Prices, PriceRecord{Time: now, Price: price})

	// Update the max and min prices
	if price > ma.Max {
		ma.Max = price
	}
	if price < ma.Min {
		ma.Min = price
	}

	// Calculate the EMA
	if len(ma.Prices) == 1 {
		ma.EMA = price
	} else {
		ma.EMA = ma.Alpha*price + (1-ma.Alpha)*ma.EMA
	}

	// Debug logging
	if math.IsNaN(ma.EMA) {
		fmt.Printf("EMA is NaN! price: %f, alpha: %f\n", price, ma.Alpha)
	}

	// Remove prices older than the interval
	for len(ma.Prices) > 0 && now.Sub(ma.Prices[0].Time) > ma.Interval {
		ma.Prices = ma.Prices[1:]

		// Recalculate max and min prices
		ma.Max = math.Inf(-1)
		ma.Min = math.Inf(1)
		for _, record := range ma.Prices {
			if record.Price > ma.Max {
				ma.Max = record.Price
			}
			if record.Price < ma.Min {
				ma.Min = record.Price
			}
		}
	}
}

func (ma *MovingAveragePrice) Average() float64 {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	return ma.EMA
}

func (ma *MovingAveragePrice) IsAboveTrend(price float64) bool {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	return price > ma.EMA
}

func (ma *MovingAveragePrice) IsBelowTrend(price float64) bool {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	return price < ma.EMA
}

func (ma *MovingAveragePrice) MaxPrice() float64 {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	return ma.Max
}

func (ma *MovingAveragePrice) MinPrice() float64 {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	return ma.Min
}
