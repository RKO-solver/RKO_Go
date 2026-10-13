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
	// Policy records the tuner's learned policy at the end of the run.
	Policy(t Table)
}

// Logger is what a Solver and its metaheuristics need from a logger.
type Logger interface {
	// SolverLogger returns the logger of a new solver.
	SolverLogger(name string) SolverLogger
	// AddSolutionPool records a new best solution in the pool.
	AddSolutionPool(cost int, elapsed time.Duration)
	// WorkerStarted records that worker id has started.
	WorkerStarted(id int, name string)
	// WorkerDone records the result of a finished worker.
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

// PolicySaver prints and saves the recorded policies.
type PolicySaver interface {
	PrintPolicies() error
	SavePoliciesCSV(prefix string) error // one file per solver: <prefix>-<solver>-<id>-policy.csv
}

// Recorder is a Logger that keeps the run's results; printer and dashboard implement it.
type Recorder interface {
	Logger
	Reporter
	CSVSaver
	PolicySaver
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

// Table is a table of text values.
type Table struct {
	Header []string
	Rows   [][]string
}
