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

`DefaultLogLevel` is `INFO`. Two helpers parse a level from a string, with different behavior for unrecognized input:

| Function | On unrecognized input | Where it's used |
|---|---|---|
| `GetLogLevel(name string) Level` | returns `INFO` | `configuration` package, YAML `logLevel` key |
| `GetLevel(level string) (Level, error)` | returns `SILENT` and an error | manual parsing where you want to reject bad input |

`GetLevelString(level Level) string` does the reverse, returning `"Silent"`, `"Info"`, or `"Verbose"`.

### Choosing a Built-in Logger

| | `channel` | `stdout` |
|---|---|---|
| Output | Live-refreshing terminal dashboard (redraws in place) | Plain lines, one per event |
| Best for | Interactive terminal runs | Batch jobs, CI, redirecting to a file |

`configuration.SolverConfiguration.LoggerType` (a `logger.LogType`, backed by the YAML `logType` key) picks between them:

```go
type LogType = uint8

const (
    CHANNEL LogType = iota // use logger/channel
    PRINT                   // use logger/stdout
)

func GetLogType(label string) LogType // "CHANNEL" or "PRINT", unknown -> PRINT
```

The `channel.Log` lifecycle methods (`Start`, `Shutdown`, `Print`, `WorkersPrint`, `GetTicker`) are driven by the solver itself — you don't call them directly. The knobs you do set yourself, before passing the logger to a solver:

- `SetTicker(timeMilliseconds int)` — dashboard refresh interval; values below 300ms are silently ignored.
- `SetNumPoolMessages(num int)` — how many recent pool entries the dashboard shows.
- `SetNumVerboseMessages(num int)` — how many recent verbose messages per solver the dashboard shows.

You can implement your own logger by satisfying the `Logger` interface.

## Usage Example

```go
import (
    "fmt"

    rko "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/logger"
    "github.com/RKO-solver/rko-go/logger/channel"
)

env := MyEnvironment{}                        // your definition.Environment
mh := []rko.MetaHeuristic{rko.GA, rko.SA}      // metaheuristics to run

log := channel.NewLoggerLevel("MyOptimizationProblem", logger.VERBOSE)

solver := rko.CreateDefaultSolver(mh, env, log)
result := solver.Solve()
fmt.Printf("Best solution: %+v\n", result)

// Retrieve and use the logged data.
reportData := log.GetReportData()
solutionData := log.GetSolutionData()
fmt.Printf("%d solver(s), %d pooled solution(s)\n", len(reportData), len(solutionData))
```

`channel.DefaultLogger("MyOptimizationProblem")` is equivalent, using `logger.DefaultLogLevel` (`INFO`) instead.

### Constructors

| Package | Function | Level | Buffer size |
|---|---|---|---|
| `channel` | `DefaultLogger(problemName string)` | `logger.DefaultLogLevel` | default |
| `channel` | `NewLoggerLevel(problemName string, level logger.Level)` | given | default |
| `channel` | `NewLogger(problemName string, logLevel logger.Level, bufferSize int)` | given | given |
| `stdout` | `DefaultLogger(problemName string)` | `logger.DefaultLogLevel` | n/a |
| `stdout` | `NewLogger(problemName string, logLevel logger.Level)` | given | n/a |

## Exporting Results to CSV

`helpers.go` provides plain CSV writers, all under package `logger`:

- `SaveCsvFile(filename string, data [][]string)` — writes rows as CSV; calls `log.Fatal` on I/O error.
- `SavePoolCsv(poolData []SolutionData, problemName string, filename ...string)` — header `cost,time`; default filename `<problemName>-pool.csv`.
- `SaveSolverCSV(solverInfo SolverInformation, problemName string, filename ...string)` — header `best,local,time`; default filename `<problemName>-<solverInfo.Name>-<solverInfo.Id>.csv`.

Both `*channel.Log` and `*stdout.Log` also have a concrete `SaveCsv(filename ...string)` method that combines the two above (one CSV per solver, plus one pool CSV). **`SaveCsv` is not part of the `logger.Logger` interface**, so if you only have a `logger.Logger` value (for example, the one returned by `configuration.CreateSolver`), you must type-assert to reach it:

```go
solver, log := configuration.CreateSolver(problemName, env, solverConfig, mhConfig)
solver.Solve()

if saver, ok := log.(interface{ SaveCsv(filename ...string) }); ok {
    saver.SaveCsv()
}
```

If you built the logger yourself with `channel.DefaultLogger` or `stdout.DefaultLogger`, you already hold the concrete type and can call `log.SaveCsv()` directly.

## Files
- `definition.go`: Contains the `Logger`, `SolverLogger` interfaces, data structures, and `LogType`.
- `constants.go`: Contains `Level` type, log level constants, and `GetLogLevel`.
- `read.go`: Contains `GetLevel` and `GetLevelString`.
- `helpers.go`: CSV export helpers (`SaveCsvFile`, `SavePoolCsv`, `SaveSolverCSV`).
- `channel/`: Channel-based logger implementation for detailed, real-time logging.
- `stdout/`: Stdout logger implementation for simple console output.

## See Also
- [Project README](../README.md)
- [definition package](../definition/README.md)

---
This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
