package channel

import (
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

type Log struct {
	updateChan chan channelMessage
	data       *information
	LogLevel   logger.Level
	ticker     time.Duration
}

func DefaultLogger() *Log {
	return NewLogger(logger.DefaultLogLevel, defaultBufferSize)
}

func NewLogger(logLevel logger.Level, bufferSize int) *Log {
	// The channel for communication
	progressChan := make(chan channelMessage, bufferSize)

	// The final data store
	store := &information{
		init:          false,
		solutionCost:  make([]int, 0, 100),
		solvers:       make([][]solverInfo, 0),
		extraMessages: make([][]string, 0),
	}

	// The logger object that workers will use
	logger := &Log{
		updateChan: progressChan,
		data:       store,
		LogLevel:   logLevel,
		ticker:     defaultTickerMilliseconds * time.Millisecond,
	}

	return logger
}

func (l *Log) AddSolutionPool(cost int) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.solutionCost = append(l.data.solutionCost, cost)
}

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
				store.registerVerbose(id, message.message)
			}

		}
	}()
}

// Shutdown closes the main channel to stop the aggregator
func (l *Log) Shutdown() {
	close(l.updateChan)
}

func (l *Log) Print() {
	l.data.printShell()
}

func (l *Log) GetTicker() time.Duration {
	return l.ticker
}

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
	l.data.solvers = append(l.data.solvers, make([]solverInfo, 0, 50))
	l.data.extraMessages = append(l.data.extraMessages, []string{})

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

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ logger.Logger = (*Log)(nil)
