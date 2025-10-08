package stdout

import (
	"fmt"

	"github.com/RKO-solver/rko-go/logger"
)

type Log struct {
	LogLevel      logger.Level
	problemName   string
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
	fmt.Println(message)
}

func (l *Log) GetLogger(name string) logger.SolverLogger {
	id := len(l.solversLogger)

	solverLogger := &SolverLogger{
		id:       id,
		name:     name,
		logLevel: l.LogLevel,
		solver:   make([]solverInfo, 0),
	}

	l.solversLogger = append(l.solversLogger, solverLogger)

	return solverLogger
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

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ logger.Logger = (*Log)(nil)
