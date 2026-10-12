package printer

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

type SolverLogger struct {
	id       int
	name     string
	level    logger.Level
	interval time.Duration
	w        io.Writer
	next     time.Duration
	lastBest int
	history  []logger.Sample // records not yet spilled to the file
	file     spill.File
	policy   *logger.Table // the tuner's learned policy; nil until Policy
}

func (s *SolverLogger) Register(local, best, iteration int, elapsed time.Duration) {
	s.history = append(s.history, logger.Sample{Local: local, Best: best, Elapsed: elapsed})
	// A full buffer is spilled synchronously, in the solver's goroutine; no lock since only this solver touches it.
	// After a spill the slice is reused from [:0]; if spilling is off or fails, it keeps growing in memory.
	if len(s.history) == cap(s.history) && s.file.Write(s.history) {
		s.history = s.history[:0]
	}

	if s.level > logger.LevelSilent && (best < s.lastBest || elapsed >= s.next) {
		s.next = elapsed + s.interval
		s.lastBest = best
		fmt.Fprintf(s.w, "(%d) %s iteration=%d local=%d best=%d time=%.3fs\n", s.id, s.name, iteration, local, best, elapsed.Seconds())
	}
}

func (s *SolverLogger) Event(a slog.Attr, elapsed time.Duration) {
	if s.level < logger.LevelVerbose {
		return
	}
	fmt.Fprintf(s.w, "(%d) %s %s time=%.3fs\n", s.id, s.name, a, elapsed.Seconds()) // a.String() is "key=value"
}

// Policy keeps t on this solver, to be read after the run.
func (s *SolverLogger) Policy(t logger.Table) {
	s.policy = &t
}

// reset deletes the spill file, applies c for this solver, empties the history and restores the initial print state.
func (s *SolverLogger) reset(c spill.Config) {
	s.next = 0
	s.lastBest = math.MaxInt
	s.policy = nil
	c.Solver, c.ID = s.name, s.id
	s.file.Reset(c)
	s.history = make([]logger.Sample, 0, spill.Capacity(c.Bytes)) // the capacity is the flush point
}

// Compile-time check that *SolverLogger implements logger.SolverLogger.
var _ logger.SolverLogger = (*SolverLogger)(nil)
