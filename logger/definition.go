// Package logger defines the Logger and SolverLogger interfaces used by RKO-Go solvers.
package logger

import (
	"io"
	"log/slog"
	"time"
)

// SolverLogger records the progress of a single running solver.
type SolverLogger interface {
	// Register records one iteration; it must not allocate.
	Register(local, best, iteration int, elapsed time.Duration)
	// Event records an occasional event, shown at LevelVerbose.
	Event(a slog.Attr, elapsed time.Duration)
}

// Logger is what a Solver and its metaheuristics need from a logger.
type Logger interface {
	// SolverLogger returns the logger of a new solver.
	SolverLogger(name string) SolverLogger
	// AddSolutionPool records a new best solution in the pool.
	AddSolutionPool(cost int, elapsed time.Duration)
	WorkerStarted(id int, name string)
	WorkerDone(r WorkerResult)
	// Reset clears the previous run; Solve calls it.
	Reset()
}

// Runner is implemented by loggers that need a loop running during Solve.
type Runner interface{ Start() (stop func()) }

// Reporter returns the recorded run.
type Reporter interface {
	Report() []SolverReport
	Pool() []PoolEntry
}

// CSVSaver saves the recorded run as CSV files.
type CSVSaver interface{ SaveCSV(prefix string) error }

// Recorder is a Logger that keeps the run's results; printer and dashboard implement it.
type Recorder interface {
	Logger
	Reporter
	CSVSaver
	io.Closer
}

// Sample is one progress record of a solver.
type Sample struct {
	Local, Best int
	Elapsed     time.Duration
}

// SolverReport is the full history of one solver.
type SolverReport struct {
	Name    string
	ID      int
	Samples []Sample
}

// PoolEntry is the cost and elapsed time of one pooled solution.
type PoolEntry struct {
	Cost    int
	Elapsed time.Duration
}

// WorkerResult is the outcome of one finished worker.
type WorkerResult struct {
	ID      int
	Name    string
	Cost    int
	Elapsed time.Duration
}
