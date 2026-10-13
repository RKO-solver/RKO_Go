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

// Log prints progress lines to a writer and keeps every solver's history.
type Log struct {
	level    logger.Level
	interval time.Duration
	w        io.Writer // not locked: the caller's writer must be safe for concurrent use
	problem  string
	hist     spill.Config // spill settings shared by every solver
	pool     []logger.PoolEntry
	solvers  []*SolverLogger
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

// SolverLogger registers a new solver and returns its logger.
func (l *Log) SolverLogger(name string) logger.SolverLogger {
	id := len(l.solvers)

	solverLogger := &SolverLogger{
		id:       id,
		name:     name,
		level:    l.level,
		interval: l.interval,
		w:        l.w,
	}
	solverLogger.reset(l.hist) // sizes the buffer; the file is created on the first spill

	l.solvers = append(l.solvers, solverLogger)

	return solverLogger
}

// Reset empties the pool and every solver's history, files and print state, keeping the registered solvers.
func (l *Log) Reset() {
	l.pool = make([]logger.PoolEntry, 0, 50)
	for _, solverLogger := range l.solvers {
		solverLogger.reset(l.hist)
	}
}

// Close deletes the last run's history files; call it after Solve returns and the results are read.
// Spilling stays off until the next Reset.
func (l *Log) Close() error {
	var err error
	for _, s := range l.solvers {
		err = errors.Join(err, s.file.Remove())
	}
	return err
}

// Report returns the history of every registered solver.
func (l *Log) Report() []logger.SolverReport {
	report := make([]logger.SolverReport, 0, len(l.solvers))

	for _, solver := range l.solvers {
		report = append(report, logger.SolverReport{
			Name:    solver.name,
			ID:      solver.id,
			Samples: solver.file.All(solver.history), // spilled records, then the in-memory tail
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
