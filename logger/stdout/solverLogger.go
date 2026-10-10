package stdout

import (
	"fmt"
	"math"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

type SolverLogger struct {
	id            int
	name          string
	logLevel      logger.Level
	printInterval float64
	lastPrint     float64
	lastBest      int
	history       []logger.Data // records not yet spilled to the file
	file          spill.File
}

func (s *SolverLogger) Register(local int, localBest int, time float64, extra string) {
	s.history = append(s.history, logger.Data{LocalCost: local, BestCost: localBest, Time: time})
	// A full buffer is spilled synchronously, in the solver's goroutine; no lock since only this solver touches it.
	// After a spill the slice is reused from [:0]; if spilling is off or fails, it keeps growing in memory.
	if len(s.history) == cap(s.history) && s.file.Write(s.history) {
		s.history = s.history[:0]
	}

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

// reset deletes the spill file, applies c for this solver, empties the history and restores the initial print state.
func (s *SolverLogger) reset(c spill.Config) {
	s.lastPrint = math.Inf(-1)
	s.lastBest = math.MaxInt
	c.Solver, c.ID = s.name, s.id
	s.file.Reset(c)
	s.history = make([]logger.Data, 0, spill.Capacity(c.Bytes)) // the capacity is the flush point
}

func (s *SolverLogger) getReportData() logger.SolverInformation {
	return logger.SolverInformation{
		Name:        s.name,
		Id:          s.id,
		Performance: s.file.All(s.history), // spilled records, then the in-memory tail
	}
}

// Compile-time check that *SolverLogger implements logger.SolverLogger.
var _ logger.SolverLogger = (*SolverLogger)(nil)
