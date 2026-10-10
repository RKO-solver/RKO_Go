// Package channel implements a live-refreshing terminal Logger for interactive runs.
package channel

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// Log is the channel-based Logger.
type Log struct {
	updateChan   chan channelMessage
	bufferSize   int
	data         *information
	LogLevel     logger.Level
	ticker       time.Duration
	problemName  string
	historyBytes int    // per-solver history buffer; 0 or less keeps everything in memory
	historyDir   string // directory of the history files; "" means the system temp directory
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

// Reset empties the previous run's data, deletes its history files and opens a fresh update channel.
func (l *Log) Reset() {
	l.data.reset(l.historyConfig())
	l.updateChan = make(chan channelMessage, l.bufferSize)
}

// Start launches the goroutine that aggregates the update channel into the run's data.
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

// Finish leaves the alternate screen, if entered, and prints the final dashboard frame on the normal screen.
func (l *Log) Finish() {
	l.data.printFinal()
}

// CleanScreen clears the dashboard while it is drawn in the alternate screen.
func (l *Log) CleanScreen() {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()
	if l.data.altScreen {
		fmt.Print(cursorHomeCode + clearToEndCode)
	}
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
	l.data.solvers = append(l.data.solvers, &solverHistory{id: id, name: name})
	l.data.solvers[id].reset(l.historyConfig()) // sizes the buffer; the file is created on the first spill
	l.data.extraMessages = append(l.data.extraMessages, make([]extraInfo, 0))

	return &SolverLoggerChannel{
		id:       id,
		name:     name,
		logLevel: l.LogLevel,
		log:      l,
	}
}

// SetHistoryBuffer sets each solver's history buffer in bytes from the next Reset on; 0 or less keeps the whole history in memory.
func (l *Log) SetHistoryBuffer(bytes int) {
	l.historyBytes = bytes
}

// SetHistoryDir sets the directory of the history files from the next Reset on; "" means the system temp directory.
func (l *Log) SetHistoryDir(dir string) {
	l.historyDir = dir
}

// historyConfig returns the spill settings shared by every solver.
func (l *Log) historyConfig() spill.Config {
	return spill.Config{Dir: l.historyDir, Bytes: l.historyBytes, Problem: l.problemName}
}

// Close deletes the last run's history files; call it after reading the results.
// Spilling stays off until the next Reset.
func (l *Log) Close() error {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	var err error
	for _, s := range l.data.solvers {
		err = errors.Join(err, s.file.Remove())
	}
	return err
}

func (l *Log) GetLogLevel() logger.Level {
	return l.LogLevel
}

// Compile-time check that *Log implements logger.Logger.
var _ logger.Logger = (*Log)(nil)
