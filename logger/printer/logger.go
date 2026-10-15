// Package printer implements a plain-text terminal Logger for batch or CI runs.
package printer

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// Log prints progress lines to a writer and keeps every metaheuristic's history.
type Log struct {
	level          logger.Level
	interval       int
	w              io.Writer // not locked: the caller's writer must be safe for concurrent use
	problem        string
	hist           spill.Config // spill settings shared by every metaheuristic
	pool           []logger.PoolEntry
	metaheuristics []*MetaheuristicLogger
}

// AddSolutionPool records a new best pool solution, printing it at LevelVerbose.
func (l *Log) AddSolutionPool(cost int, elapsed time.Duration) {
	l.pool = append(l.pool, logger.PoolEntry{Cost: cost, Elapsed: elapsed})

	if l.level >= logger.LevelVerbose {
		fmt.Fprintf(l.w, "pool cost=%d time=%.3fs\n", cost, elapsed.Seconds())
	}
}

// WorkerStarted prints that worker id has started, unless the level is LevelSilent.
func (l *Log) WorkerStarted(id int, name string) {
	if l.level > logger.LevelSilent {
		fmt.Fprintf(l.w, "(%d) %s started\n", id, name)
	}
}

// WorkerDone prints the result of a finished worker, unless the level is LevelSilent.
func (l *Log) WorkerDone(r logger.WorkerResult) {
	if l.level > logger.LevelSilent {
		fmt.Fprintf(l.w, "(%d) %s done cost=%d time=%.3fs\n", r.ID, r.Name, r.Cost, r.Elapsed.Seconds())
	}
}

// MetaheuristicLogger registers a new metaheuristic and returns its logger.
func (l *Log) MetaheuristicLogger(name string) logger.MetaheuristicLogger {
	id := len(l.metaheuristics)

	metaheuristicLogger := &MetaheuristicLogger{
		id:       id,
		name:     name,
		level:    l.level,
		interval: max(l.interval, 0),
		w:        l.w,
	}
	metaheuristicLogger.reset(l.hist) // sizes the buffer; the file is created on the first spill

	l.metaheuristics = append(l.metaheuristics, metaheuristicLogger)

	return metaheuristicLogger
}

// Reset empties the pool and every metaheuristic's history, files and print state, keeping the registered metaheuristics.
func (l *Log) Reset() {
	l.pool = make([]logger.PoolEntry, 0, 50)
	for _, metaheuristicLogger := range l.metaheuristics {
		metaheuristicLogger.reset(l.hist)
	}
}

// Close deletes the last run's history files; call it after Solve returns and the results are read.
// Spilling stays off until the next Reset.
func (l *Log) Close() error {
	var err error
	for _, s := range l.metaheuristics {
		err = errors.Join(err, s.file.Remove())
	}
	return err
}

// Report returns the history of every registered metaheuristic.
func (l *Log) Report() []logger.MetaheuristicReport {
	report := make([]logger.MetaheuristicReport, 0, len(l.metaheuristics))

	for _, metaheuristic := range l.metaheuristics {
		report = append(report, logger.MetaheuristicReport{
			Name:    metaheuristic.name,
			ID:      metaheuristic.id,
			Samples: metaheuristic.file.All(metaheuristic.history), // spilled records, then the in-memory tail
		})
	}

	return report
}

// Pool returns a copy of the recorded pool solutions.
func (l *Log) Pool() []logger.PoolEntry {
	return slices.Clone(l.pool)
}

// Compile-time check that *Log implements logger.Recorder.
var _ logger.Recorder = (*Log)(nil)
