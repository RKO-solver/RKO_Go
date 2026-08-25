// Package channel implements a live-refreshing terminal Logger for interactive runs.
package channel

import (
	"fmt"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

// Log is the channel-based Logger.
type Log struct {
	updateChan  chan channelMessage
	data        *information
	LogLevel    logger.Level
	ticker      time.Duration
	problemName string
}

func (l *Log) AddSolutionPool(cost int, time float64) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.pool = append(l.data.pool, poolInfo{cost: cost, time: time})
}

func (l *Log) WorkerDone(message string) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.workerMessages = append(l.data.workerMessages, message)
}

// Start launches the goroutine that aggregates channel messages.
func (l *Log) Start(aggregatorWg *sync.WaitGroup) {
	progressChan := l.updateChan
	store := l.data

	// Start the single aggregator goroutine in the background.
	// It's the only one that writes to 'store'.
	aggregatorWg.Add(1)
	go func() {
		defer aggregatorWg.Done()
		for message := range progressChan {
			id := message.id
			switch message.t {
			case infoMessage:
				store.registerInfo(id, message.info)
			case verboseMessage:
				store.registerVerbose(id, message.extra)
			}

		}
	}()
}

// SetNumPoolMessages sets how many pool messages the dashboard shows.
func (l *Log) SetNumPoolMessages(num int) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()
	l.data.numLinesPool = num
}

// SetNumVerboseMessages sets how many verbose messages the dashboard shows.
func (l *Log) SetNumVerboseMessages(num int) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()
	l.data.numVerboseMessages = num
}

// Shutdown closes the update channel to stop the aggregator.
func (l *Log) Shutdown() {
	close(l.updateChan)
}

// Print redraws the live dashboard.
func (l *Log) Print() {
	l.data.printShell()
}

// CleanScreen erases the previously printed dashboard lines.
func (l *Log) CleanScreen() {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()
	fmt.Print(fmt.Sprintf("\033[%dA", l.data.previousLineCount))
	for _ = range l.data.previousLineCount {
		fmt.Print(cleanLineCode)
	}
	fmt.Print(fmt.Sprintf("\033[%dA", l.data.previousLineCount))
}

// WorkersPrint prints all recorded worker-done messages.
func (l *Log) WorkersPrint() {
	for _, message := range l.data.workerMessages {
		fmt.Println(message)
	}
}

// GetTicker returns the dashboard refresh interval.
func (l *Log) GetTicker() time.Duration {
	return l.ticker
}

// SetTicker sets the dashboard refresh interval, silently ignoring values below 300ms.
func (l *Log) SetTicker(timeMilliseconds int) {
	if timeMilliseconds < minimumTickerMilliseconds {
		return
	}
	l.ticker = time.Duration(timeMilliseconds) * time.Millisecond
}

func (l *Log) GetLogger(name string) logger.SolverLogger {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	id := len(l.data.solvers)
	l.data.solvers = append(l.data.solvers, make([]solverInfo, 0))
	l.data.extraMessages = append(l.data.extraMessages, make([]extraInfo, 0))

	return &SolverLoggerChannel{
		id:         id,
		name:       name,
		logLevel:   l.LogLevel,
		updateChan: l.updateChan,
	}
}

func (l *Log) GetLogLevel() logger.Level {
	return l.LogLevel
}

// Compile-time check that *Log implements logger.Logger.
var _ logger.Logger = (*Log)(nil)
