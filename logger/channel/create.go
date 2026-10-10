package channel

import (
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

// DefaultLogger creates a Log with the default level and buffer size.
func DefaultLogger(problemName string) *Log {
	return NewLogger(problemName, logger.DefaultLogLevel, defaultBufferSize)
}

// NewLoggerLevel creates a Log at the given level.
func NewLoggerLevel(problemName string, level logger.Level) *Log {
	return NewLogger(problemName, level, defaultBufferSize)
}

// NewLogger creates a Log with the given level and channel buffer size.
func NewLogger(problemName string, logLevel logger.Level, bufferSize int) *Log {
	// The final data store
	store := &information{
		numLinesPool:       numLinesSolutionPool,
		numVerboseMessages: lastNVerboses,
		pool:               make([]poolInfo, 0, 100),
		solvers:            make([]*solverHistory, 0),
		extraMessages:      make([][]extraInfo, 0),
		workerMessages:     make([]string, 0),
	}

	return &Log{
		bufferSize:   bufferSize,
		data:         store,
		LogLevel:     logLevel,
		ticker:       defaultTickerMilliseconds * time.Millisecond,
		problemName:  problemName,
		historyBytes: defaultHistoryBytes,
	}
}
