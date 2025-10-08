package channel

import (
	"time"

	"github.com/RKO-solver/rko-go/logger"
)

func DefaultLogger(problemName string) *Log {
	return NewLogger(problemName, logger.DefaultLogLevel, defaultBufferSize)
}

func NewLoggerLevel(problemName string, level logger.Level) *Log {
	return NewLogger(problemName, level, defaultBufferSize)
}

func NewLogger(problemName string, logLevel logger.Level, bufferSize int) *Log {
	// The channel for communication
	progressChan := make(chan channelMessage, bufferSize)

	// The final data store
	store := &information{
		init:              false,
		previousLineCount: 0,
		pool:              make([]poolInfo, 0, 100),
		solvers:           make([][]solverInfo, 0),
		extraMessages:     make([][]extraInfo, 0),
		workerMessages:    make([]string, 0),
	}

	return &Log{
		updateChan:  progressChan,
		data:        store,
		LogLevel:    logLevel,
		ticker:      defaultTickerMilliseconds * time.Millisecond,
		problemName: problemName,
	}
}
