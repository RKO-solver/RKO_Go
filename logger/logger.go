package logger

import (
	"sync"
	"time"
)

type Logger struct {
	updateChan chan channelMessage
	data       *information
	LogLevel   Level
	ticker     time.Duration
}

func DefaultLogger() *Logger {
	return NewLogger(defaultLogLevel, defaultBufferSize)
}

func NewLogger(logLevel Level, bufferSize int) *Logger {
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
	logger := &Logger{
		updateChan: progressChan,
		data:       store,
		LogLevel:   logLevel,
		ticker:     defaultTickerMilliseconds * time.Millisecond,
	}

	return logger
}

func (l *Logger) AddSolutionPool(cost int) {
	l.data.mu.Lock()
	defer l.data.mu.Unlock()

	l.data.solutionCost = append(l.data.solutionCost, cost)
}

func (l *Logger) Start(aggregatorWg *sync.WaitGroup) {
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
func (l *Logger) Shutdown() {
	close(l.updateChan)
}

func (l *Logger) Print() {
	l.data.printShell()
}

func (l *Logger) GetTicker() time.Duration {
	return l.ticker
}

func (l *Logger) SetTicker(timeMilliseconds int) {
	if timeMilliseconds < minimumTickerMilliseconds {
		return
	}
	l.ticker = time.Duration(timeMilliseconds) * time.Millisecond
}
