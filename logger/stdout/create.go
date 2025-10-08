package stdout

import (
	"github.com/RKO-solver/rko-go/logger"
)

func DefaultLogger(problemName string) *Log {
	return NewLogger(problemName, logger.DefaultLogLevel)
}

func NewLogger(problemName string, logLevel logger.Level) *Log {
	return &Log{
		LogLevel:      logLevel,
		problemName:   problemName,
		pool:          make([]poolInfo, 0),
		solversLogger: make([]*SolverLogger, 0),
	}
}
