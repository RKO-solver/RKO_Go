// Package stdout implements a plain-text terminal Logger for batch or CI runs.
package stdout

import (
	"fmt"

	"github.com/RKO-solver/rko-go/logger"
)

type Log struct {
	LogLevel      logger.Level
	problemName   string
	printInterval float64
	pool          []poolInfo
	solversLogger []*SolverLogger
}

func (l *Log) AddSolutionPool(cost int, time float64) {
	l.pool = append(l.pool, poolInfo{cost: cost, time: time})

	if l.LogLevel >= logger.VERBOSE {
		fmt.Printf("Adding to pool %d at %.3fs\n", cost, time)
	}
}

func (l *Log) WorkerDone(message string) {
	if l.LogLevel > logger.SILENT {
		fmt.Println(message)
	}
}

func (l *Log) GetLogger(name string) logger.SolverLogger {
	id := len(l.solversLogger)

	solverLogger := &SolverLogger{
		id:            id,
		name:          name,
		logLevel:      l.LogLevel,
		printInterval: l.printInterval,
	}
	solverLogger.reset()

	l.solversLogger = append(l.solversLogger, solverLogger)

	return solverLogger
}

// Reset empties the pool and every solver's history and print state, keeping the registered solvers.
func (l *Log) Reset() {
	l.pool = make([]poolInfo, 0)
	for _, solverLogger := range l.solversLogger {
		solverLogger.reset()
	}
}

// SetPrintInterval sets the minimum time between two progress prints of the same solver.
func (l *Log) SetPrintInterval(timeMilliseconds int) {
	l.printInterval = float64(timeMilliseconds) / 1000
	for _, solverLogger := range l.solversLogger {
		solverLogger.printInterval = l.printInterval
	}
}

func (l *Log) GetLogLevel() logger.Level {
	return l.LogLevel
}

func (l *Log) GetReportData() []logger.SolverInformation {
	report := make([]logger.SolverInformation, 0, len(l.solversLogger))

	for _, solver := range l.solversLogger {
		solverReport := solver.getReportData()

		report = append(report, solverReport)
	}

	return report
}

func (l *Log) GetSolutionData() []logger.SolutionData {
	report := make([]logger.SolutionData, 0, len(l.pool))

	for _, pool := range l.pool {
		report = append(report, logger.SolutionData{
			Cost: pool.cost,
			Time: pool.time,
		})
	}

	return report
}

// Compile-time check that *Log implements logger.Logger.
var _ logger.Logger = (*Log)(nil)
