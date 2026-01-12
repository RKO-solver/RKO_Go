# `logger` Package

This package provides logging interfaces and implementations for tracking the progress and results of metaheuristic optimization in the RKO-Go library.

## Overview

Logging is essential for monitoring the optimization process, debugging, and saving results. The `logger` package defines flexible interfaces and includes built-in loggers for different use cases.

## Main Components

### Interfaces

The core of this package consists of two interfaces:

**Logger Interface** - Implement this to create custom loggers:

```go
type Logger interface {
    AddSolutionPool(cost int, time float64)
    WorkerDone(message string)
    GetLogger(name string) SolverLogger
    GetLogLevel() Level
    GetReportData() []SolverInformation
    GetSolutionData() []SolutionData
}
```
- `AddSolutionPool`: Report a new best solution found across all solvers.
- `WorkerDone`: Called when a solver finishes execution.
- `GetLogger`: Returns a `SolverLogger` for a specific solver by name.
- `GetLogLevel`: Returns the current logging level (`SILENT`, `INFO`, or `VERBOSE`).
- `GetReportData`: Returns performance data for all solvers.
- `GetSolutionData`: Returns the best solution data found.

**SolverLogger Interface** - Implement this for solver-specific logging:

```go
type SolverLogger interface {
    Register(local int, localBest int, time float64, extra string)
    Verbose(message string, timeStamp float64)
}
```
- `Register`: Record progress from a solver (local cost, best cost, elapsed time, and optional extra data).
- `Verbose`: Log detailed messages with timestamps.

### Data Structures

- `SolverInformation`: Contains solver name, ID, and performance data.
- `Data`: Represents a single performance record (local cost, best cost, elapsed time).
- `SolutionData`: Stores best solution cost and time found.

### Built-in Loggers

- **Channel Logger** (`channel/`): Provides detailed logging with communication channels for real-time progress tracking. Use `channel.DefaultLogger(problemName)` or `channel.NewLoggerLevel(problemName, level)` to instantiate.
- **Stdout Logger** (`stdout/`): Prints progress and messages directly to the standard output. Use `stdout.DefaultLogger(problemName)` or `stdout.NewLogger(problemName, level)` to instantiate.

### Log Levels

The logger supports three logging levels:

```go
type Level uint8

const (
    SILENT   Level = iota  // No output
    INFO                    // Basic progress information (default)
    VERBOSE                 // Detailed progress information
)
```

You can implement your own logger by satisfying the `Logger` interface.

## Usage Example

```go
import (
    "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/logger/channel"
)

// Create a logger with default INFO level
logger := channel.DefaultLogger("MyOptimizationProblem")

// Or create with a specific log level
logger := channel.NewLoggerLevel("MyOptimizationProblem", logger.VERBOSE)
```

Pass your logger to the solver when creating it:

```go
solver := rko.CreateDefaultSolver(mh, env, logger)
result := solver.Solve()

// Retrieve and use the logged data
reportData := logger.GetReportData()
solutionData := logger.GetSolutionData()
```

## Files
- `definition.go`: Contains the `Logger`, `SolverLogger` interfaces and data structures.
- `constants.go`: Contains `Level` type and log level constants.
- `helpers.go`: Helper functions for logging utilities.
- `channel/`: Channel-based logger implementation for detailed, real-time logging.
- `stdout/`: Stdout logger implementation for simple console output.

## See Also
- [Project README](../README.md)
- [definition package](../definition/README.md)

---
This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
