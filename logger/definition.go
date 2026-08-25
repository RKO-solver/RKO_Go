// Package logger defines the Logger and SolverLogger interfaces used by RKO-Go solvers.
package logger

import "strings"

// SolverInformation holds one solver's name, ID, and full performance history.
type SolverInformation struct {
	Name        string
	Id          int
	Performance []Data
}

// SolutionData holds the cost and elapsed time of one pooled solution.
type SolutionData struct {
	Cost int
	Time float64
}

// Data holds one progress sample recorded by a SolverLogger.
type Data struct {
	LocalCost int
	BestCost  int
	Time      float64
}

// Logger aggregates solution-pool updates and per-solver progress across a run.
type Logger interface {
	AddSolutionPool(cost int, time float64)
	WorkerDone(message string)
	GetLogger(name string) SolverLogger
	GetLogLevel() Level
	GetReportData() []SolverInformation
	GetSolutionData() []SolutionData
}

// SolverLogger records progress messages for a single running solver.
type SolverLogger interface {
	Register(local int, localBest int, time float64, extra string)
	Verbose(message string, timeStamp float64)
}

// LogType selects which Logger implementation a solver uses.
type LogType = uint8

const (
	CHANNEL LogType = iota
	PRINT
)

// GetLogType parses a config label into a LogType, defaulting to PRINT when unrecognized.
func GetLogType(label string) LogType {
	label = strings.ToUpper(label)
	switch label {
	case "CHANNEL":
		return CHANNEL
	case "PRINT":
		return PRINT
	default:
		return PRINT
	}
}
