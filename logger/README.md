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
    Reset()
    Close() error
}
```
- `AddSolutionPool`: Report a new best solution found across all solvers.
- `WorkerDone`: Called when a solver finishes execution.
- `GetLogger`: Returns a `SolverLogger` for a specific solver by name.
- `GetLogLevel`: Returns the current logging level (`SILENT`, `INFO`, or `VERBOSE`).
- `GetReportData`: Returns performance data for all solvers.
- `GetSolutionData`: Returns the best solution data found.
- `Reset`: Called by `Solve` at the start of every run; empties the previous run's data.
- `Close`: Deletes the last run's history files; safe to call twice. `Solve` never calls it, so add `defer log.Close()` after creating the logger. Read or save the results before `Close`: afterwards history stays in memory, and spilling to disk resumes at the next `Solve`.

**SolverLogger Interface** - Implement this for solver-specific logging:

```go
type SolverLogger interface {
    Register(local int, localBest int, time float64, extra string)
    Verbose(message string, timeStamp float64)
}
```
- `Register`: Record progress from a solver (local cost, best cost, elapsed time, and optional extra data).
- `Verbose`: Log an event, such as a population restart, with its timestamp. Per-iteration progress goes in `Register`'s `extra`.

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
| Output | Live-refreshing terminal dashboard (redraws in place) | Plain lines; solver progress on improvement or once per print interval |
| Best for | Interactive terminal runs | Batch jobs, CI, redirecting to a file |
| History buffer per solver | 128 KB | 512 KB |

`configuration.SolverConfiguration.LoggerType` (a `logger.LogType`, backed by the YAML `logType` key) picks between them:

```go
type LogType = uint8

const (
    CHANNEL LogType = iota // use logger/channel
    PRINT                   // use logger/stdout
)

func GetLogType(label string) LogType // "CHANNEL" or "PRINT", unknown -> PRINT
```

The `channel` dashboard draws in the terminal's alternate screen. When the run ends, it leaves the alternate screen and prints the final frame, followed by each solver's done message, to the normal terminal, so the result stays in the scrollback.

The `channel.Log` lifecycle methods (`Start`, `Shutdown`, `Print`, `Finish`, `WorkersPrint`, `GetTicker`) are driven by the solver itself — you don't call them directly. The knobs you do set yourself, before passing the logger to a solver:

- `SetTicker(timeMilliseconds int)` — dashboard refresh interval; values below 300ms are silently ignored.
- `SetNumPoolMessages(num int)` — how many recent pool entries the dashboard shows.
- `SetNumVerboseMessages(num int)` — how many recent verbose messages per solver the dashboard shows.

The `stdout.Log` prints a solver's progress when its best cost improves, otherwise at most once per print interval (default 500 ms). Verbose messages are always printed.

- `SetPrintInterval(timeMilliseconds int)` — minimum time between two prints of the same solver.

With both loggers, each `Solve` empties the previous run's information and deletes its history files, so `GetReportData` and `GetSolutionData` return the last run.

### Solver History on Disk

Both loggers record every `Register` call. Each solver keeps its records in a memory buffer, 24 bytes per record (128 KB per solver for `channel`, 512 KB for `stdout`). When the buffer is full, the logger appends it to that solver's history file in one write and reuses the buffer. After its first flush, a solver also keeps an encoding buffer of the same size. The file is created on the first flush, so a run that never fills the buffer never touches the disk.

`GetReportData` and `SaveCsv` return the whole history: the file's records first, then the ones still in memory. `SaveCsv` streams the file into the CSV; `GetReportData` returns the whole history in one slice.

| Method | Effect |
|---|---|
| `SetHistoryBuffer(bytes int)` | Buffer size per solver. `0` keeps the whole history in memory. |
| `SetHistoryDir(dir string)` | Directory of the history files. Default: the system temp directory (`os.TempDir()`). |
| `Close() error` | Part of `Logger`. Deletes the last run's history files. Call it after reading the results. |

Both setters take effect at the next `Solve`, so you can call them after creating the solver. Each `Solve` deletes the previous run's history files. If a solver's file cannot be created or written, that solver keeps the rest of its history in memory and the logger prints one warning to stderr.

Each file is named `rko-history-<problem>-<solver>-<id>-*.bin`, with characters invalid in file names replaced by `_` and each name part cut to 32 bytes. It holds fixed 24-byte little-endian records: local cost (`int64`), best cost (`int64`), time in seconds (`float64`).

The setters are not part of the `logger.Logger` interface. With the value returned by `configuration.CreateSolver`, type-assert to reach them:

```go
solver, log := configuration.CreateSolver(problemName, env, solverConfig)
defer log.Close()
if l, ok := log.(interface{ SetHistoryBuffer(bytes int) }); ok {
    l.SetHistoryBuffer(1 << 20)
}
solver.Solve(context.Background())
// read GetReportData or call SaveCsv here
```

You can implement your own logger by satisfying the `Logger` interface, including `Close() error`; return `nil` if there is nothing to release.

## Usage Example

```go
import (
    "context"
    "fmt"

    rko "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/logger"
    "github.com/RKO-solver/rko-go/logger/channel"
)

env := MyEnvironment{}                        // your definition.Environment
mh := []rko.MetaHeuristic{rko.GA, rko.SA}      // metaheuristics to run

log := channel.NewLoggerLevel("MyOptimizationProblem", logger.VERBOSE)
defer log.Close() // runs after the results below are read

solver := rko.CreateDefaultSolver(mh, env, log)
result := solver.Solve(context.Background())
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
solver, log := configuration.CreateSolver(problemName, env, solverConfig)
defer log.Close()
solver.Solve(context.Background())

if saver, ok := log.(interface{ SaveCsv(filename ...string) }); ok {
    saver.SaveCsv()
}
```

If you built the logger yourself with `channel.DefaultLogger` or `stdout.DefaultLogger`, you already hold the concrete type and can call `log.SaveCsv()` directly. `SaveCsv` reads each solver's history file as it writes the CSV, so it never loads the whole history in memory.

## Files
- `definition.go`: Contains the `Logger`, `SolverLogger` interfaces, data structures, and `LogType`.
- `constants.go`: Contains `Level` type, log level constants, and `GetLogLevel`.
- `read.go`: Contains `GetLevel` and `GetLevelString`.
- `helpers.go`: CSV export helpers (`SaveCsvFile`, `SavePoolCsv`, `SaveSolverCSV`).
- `internal/spill/`: Writes each solver's history to a temporary binary file and reads it back, for both built-in loggers.
- `channel/`: Channel-based logger implementation for detailed, real-time logging.
- `stdout/`: Stdout logger implementation for simple console output.

## See Also
- [Project README](../README.md)
- [definition package](../definition/README.md)

---
This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
