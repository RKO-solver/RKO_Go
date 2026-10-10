// Package stdout implements a plain-text terminal Logger for batch or CI runs.
package stdout

import (
	"errors"
	"fmt"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

type Log struct {
	LogLevel      logger.Level
	problemName   string
	printInterval float64
	historyBytes  int    // per-solver history buffer; 0 or less keeps everything in memory
	historyDir    string // directory of the history files; "" means the system temp directory
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
	solverLogger.reset(l.historyConfig()) // sizes the buffer; the file is created on the first spill

	l.solversLogger = append(l.solversLogger, solverLogger)

	return solverLogger
}

// Reset empties the pool and every solver's history, files and print state, keeping the registered solvers.
func (l *Log) Reset() {
	l.pool = make([]poolInfo, 0, 50)
	for _, solverLogger := range l.solversLogger {
		solverLogger.reset(l.historyConfig())
	}
}

// SetPrintInterval sets the minimum time between two progress prints of the same solver.
func (l *Log) SetPrintInterval(timeMilliseconds int) {
	l.printInterval = float64(timeMilliseconds) / 1000
	for _, solverLogger := range l.solversLogger {
		solverLogger.printInterval = l.printInterval
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

// Close deletes the last run's history files; call it after Solve returns and the results are read.
// Spilling stays off until the next Reset.
func (l *Log) Close() error {
	var err error
	for _, s := range l.solversLogger {
		err = errors.Join(err, s.file.Remove())
	}
	return err
}

func (l *Log) GetLogLevel() logger.Level {
	return l.LogLevel
}

func (l *Log) GetReportData() []logger.SolverInformation {
	report := make([]logger.SolverInformation, 0, len(l.solversLogger))

	for _, solver := range l.solversLogger {
		report = append(report, solver.getReportData())
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
