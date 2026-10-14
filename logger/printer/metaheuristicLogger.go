package printer

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// MetaheuristicLogger prints and records the progress of one metaheuristic.
type MetaheuristicLogger struct {
	id       int
	name     string
	level    logger.Level
	interval int
	count    int
	w        io.Writer
	history  []logger.Sample // records not yet spilled to the file
	file     spill.File
	policy   *logger.Table // the tuner's learned policy; nil until Policy
}

func (s *MetaheuristicLogger) Register(local, best, iteration int, elapsed time.Duration) {
	s.history = append(s.history, logger.Sample{Local: local, Best: best, Elapsed: elapsed})
	// A full buffer is spilled synchronously, in the metaheuristic's goroutine; no lock since only this metaheuristic touches it.
	// After a spill the slice is reused from [:0]; if spilling is off or fails, it keeps growing in memory.
	if len(s.history) == cap(s.history) && s.file.Write(s.history) {
		s.history = s.history[:0]
	}

	if s.level > logger.LevelSilent && s.count <= 0 {
		s.count = s.interval + 1
		fmt.Fprintf(s.w, "(%d) %s iteration=%d local=%d best=%d time=%.3fs\n", s.id, s.name, iteration, local, best, elapsed.Seconds())
	}
	s.count--
}

func (s *MetaheuristicLogger) Event(a slog.Attr, elapsed time.Duration) {
	if s.level < logger.LevelVerbose {
		return
	}
	fmt.Fprintf(s.w, "(%d) %s %s time=%.3fs\n", s.id, s.name, a, elapsed.Seconds()) // a.String() is "key=value"
}

// Policy keeps t on this metaheuristic, to be read after the run.
func (s *MetaheuristicLogger) Policy(t logger.Table) {
	s.policy = &t
}

// reset deletes the spill file, applies c for this metaheuristic, empties the history and restores the initial print state.
func (s *MetaheuristicLogger) reset(c spill.Config) {
	s.count = 0
	s.policy = nil
	c.Metaheuristic, c.ID = s.name, s.id
	s.file.Reset(c)
	s.history = make([]logger.Sample, 0, spill.Capacity(c.Bytes)) // the capacity is the flush point
}

// Compile-time check that *MetaheuristicLogger implements logger.MetaheuristicLogger.
var _ logger.MetaheuristicLogger = (*MetaheuristicLogger)(nil)
