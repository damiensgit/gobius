package metrics

import (
	"context"
	"fmt"
	"math"
	"os"
	"sync/atomic"
	"time"

	"github.com/olekukonko/tablewriter"
)

type TaskTracker struct {
	successCountTotal int64 // for entire session
	failCountTotal    int64 // for entire session
	successCount      int64
	failCount         int64
	totalRatio        float64
	measurements      int64
	solvedCount       int64 // track gpu solve count
	solvedCountTotal  int64 // track gpu solve count total
	sessionStart      time.Time
	silence           bool
}

func NewTaskTracker(appContext context.Context) *TaskTracker {
	t := &TaskTracker{}
	return t
}

// Add a new method to start the logging goroutine
func (t *TaskTracker) StartLogging(appContext context.Context) {
	t.sessionStart = time.Now()
	go t.logEveryMinute(appContext)
}

func (t *TaskTracker) Silence(silence bool) {
	t.silence = silence
}

func (t *TaskTracker) logEveryMinute(appContext context.Context) {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-appContext.Done():
			ticker.Stop()
			return
		case <-ticker.C:

			if t.silence {
				continue
			}
			sc := atomic.LoadInt64(&t.successCount)
			fc := atomic.LoadInt64(&t.failCount)
			sct := atomic.LoadInt64(&t.successCountTotal)
			fct := atomic.LoadInt64(&t.failCountTotal)

			solvecount := atomic.LoadInt64(&t.solvedCount)
			solvecounttotal := atomic.LoadInt64(&t.solvedCountTotal)

			total := sc + fc
			totalSession := sct + fct
			ratio := math.NaN()
			averageRatio := math.NaN()
			averageTasksPerPeriod := math.NaN()
			averageSolvesPerPeriod := math.NaN()
			if total > 0 {
				measurements := atomic.AddInt64(&t.measurements, 1)
				ratio = float64(sc) / float64(total)
				t.totalRatio += ratio
				averageRatio = t.totalRatio / float64(measurements)
				averageTasksPerPeriod = float64(totalSession) / float64(measurements)
				averageSolvesPerPeriod = float64(solvecounttotal) / float64(measurements)
			}

			sessionDuration := time.Since(t.sessionStart)
			hours := int(sessionDuration.Hours())
			minutes := int(sessionDuration.Minutes()) % 60
			seconds := int(sessionDuration.Seconds()) % 60
			formattedDuration := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

			data := [][]string{}

			data = append(data, []string{"Session Time", formattedDuration, "", ""})
			data = append(data, []string{"Past 1m Solved", fmt.Sprintf("%d", solvecount), "", ""})
			data = append(data, []string{"Past 1m Solutions", "", fmt.Sprintf("%d/%d (%.2f%%)", sc, total, ratio*100), fmt.Sprintf("%d/%d (%.2f%%)", fc, total, (1-ratio)*100)})
			data = append(data, []string{"Average Solution Success", fmt.Sprintf("%.2f%%", averageRatio*100), "", ""})
			data = append(data, []string{"Average Solutions/Minute", fmt.Sprintf("%.2f", averageTasksPerPeriod), "", ""})
			data = append(data, []string{"Average Solves/Minute", fmt.Sprintf("%.2f", averageSolvesPerPeriod), "", ""})

			sessionRatio := math.NaN()
			if totalSession > 0 {
				sessionRatio = float64(sct) / float64(totalSession)
			}

			data = append(data, []string{"Solved", fmt.Sprintf("%d", solvecounttotal), "", ""})
			data = append(data, []string{"Solution Success", "", fmt.Sprintf("%d/%d (%.2f%%)", sct, totalSession, sessionRatio*100), fmt.Sprintf("%d/%d (%.2f%%)", fct, totalSession, (1-sessionRatio)*100)})

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"SESSION METRICS", "VALUE", "SUCCESS", "FAIL"})
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

			table.SetAlignment(tablewriter.ALIGN_RIGHT)
			table.AppendBulk(data)
			table.Render()

			atomic.StoreInt64(&t.successCount, 0)
			atomic.StoreInt64(&t.failCount, 0)
			atomic.StoreInt64(&t.solvedCount, 0)
		}
	}
}

func (t *TaskTracker) AverageTasksPerPeriod() float64 {
	sct := atomic.LoadInt64(&t.successCountTotal)
	fct := atomic.LoadInt64(&t.failCountTotal)
	measurements := atomic.LoadInt64(&t.measurements)

	totalSession := sct + fct

	return float64(totalSession) / float64(measurements)
}

func (t *TaskTracker) TaskSucceeded() {
	atomic.AddInt64(&t.successCount, 1)
	atomic.AddInt64(&t.successCountTotal, 1)
}

func (t *TaskTracker) TaskFailed() {
	atomic.AddInt64(&t.failCount, 1)
	atomic.AddInt64(&t.failCountTotal, 1)
}

func (t *TaskTracker) Solved() {
	atomic.AddInt64(&t.solvedCount, 1)
	atomic.AddInt64(&t.solvedCountTotal, 1)
}

func (t *TaskTracker) GetSolvedLastMinute() int64 {
	return atomic.LoadInt64(&t.solvedCount)
}

func (t *TaskTracker) GetSuccessCount() int64 {
	return atomic.LoadInt64(&t.successCount)
}

func (t *TaskTracker) GetTotalCount() int64 {
	return atomic.LoadInt64(&t.successCount) + atomic.LoadInt64(&t.failCount)
}

func (t *TaskTracker) GetSuccessRate() float64 {
	sc := atomic.LoadInt64(&t.successCount)
	fc := atomic.LoadInt64(&t.failCount)
	total := sc + fc
	if total > 0 {
		return float64(sc) / float64(total)
	}
	return 0
}

func (t *TaskTracker) GetAverageSolutionRate() float64 {
	measurements := atomic.LoadInt64(&t.measurements)
	if measurements > 0 {
		return t.totalRatio / float64(measurements)
	}
	return 0
}

func (t *TaskTracker) GetAverageSolutionsPerMin() float64 {
	sct := atomic.LoadInt64(&t.successCountTotal)
	fct := atomic.LoadInt64(&t.failCountTotal)
	measurements := atomic.LoadInt64(&t.measurements)
	if measurements > 0 {
		return float64(sct+fct) / float64(measurements)
	}
	return 0
}

func (t *TaskTracker) GetAverageSolvesPerMin() float64 {
	solvecounttotal := atomic.LoadInt64(&t.solvedCountTotal)
	measurements := atomic.LoadInt64(&t.measurements)
	if measurements > 0 {
		return float64(solvecounttotal) / float64(measurements)
	}
	return 0
}

func (t *TaskTracker) GetAverageTasksPerPeriod() float64 {
	sct := atomic.LoadInt64(&t.successCountTotal)
	fct := atomic.LoadInt64(&t.failCountTotal)
	measurements := atomic.LoadInt64(&t.measurements)
	if measurements > 0 {
		return float64(sct+fct) / float64(measurements)
	}
	return 0
}
