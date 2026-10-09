package stdout

import (
	"fmt"
	"math"

	"github.com/RKO-solver/rko-go/logger"
)

type SolverLogger struct {
	id            int
	name          string
	logLevel      logger.Level
	printInterval float64
	lastPrint     float64
	lastBest      int
	solver        []solverInfo
}

func (s *SolverLogger) Register(local int, localBest int, time float64, extra string) {
	s.solver = append(s.solver, solverInfo{
		localBest: localBest,
		local:     local,
		time:      time,
	})

	if s.logLevel > logger.SILENT && (localBest < s.lastBest || time-s.lastPrint >= s.printInterval) {
		s.lastPrint = time
		s.lastBest = localBest
		fmt.Printf("(%d) %s: %s\n\tLocal %d Best %d\n\tTime %.3fs\n", s.id, s.name, extra, local, localBest, time)
	}
}

func (s *SolverLogger) Verbose(message string, timeStamp float64) {
	if s.logLevel >= logger.VERBOSE {
		fmt.Printf("(%d) %s Verbose: %s Time %.3fs\n", s.id, s.name, message, timeStamp)
	}
}

// reset empties the history and restores the initial print state.
func (s *SolverLogger) reset() {
	s.lastPrint = math.Inf(-1)
	s.lastBest = math.MaxInt
	s.solver = make([]solverInfo, 0)
}

func (s *SolverLogger) getReportData() logger.SolverInformation {
	solverReport := logger.SolverInformation{
		Name:        s.name,
		Id:          s.id,
		Performance: make([]logger.Data, 0),
	}

	for _, info := range s.solver {
		solverReport.Performance = append(solverReport.Performance, logger.Data{
			LocalCost: info.local,
			BestCost:  info.localBest,
			Time:      info.time,
		})
	}

	return solverReport
}

// Compile-time check that *SolverLogger implements logger.SolverLogger.
var _ logger.SolverLogger = (*SolverLogger)(nil)
