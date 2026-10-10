# `logger` Package

The loggers RKO-Go solvers report to. A logger prints progress, records each solver's history and every solution-pool improvement, and saves them to CSV. Two ship with the library; implement `Logger` to write your own.

## Quick start

```go
log := channel.NewLoggerLevel("MyProblem", logger.VERBOSE)
defer log.Close() // deletes the history files once the results below are read

solver := rko.CreateDefaultSolver([]rko.MetaHeuristic{rko.GA, rko.SA}, env, log)
result := solver.Solve(context.Background())
fmt.Printf("Best solution: %+v\n", result)

report := log.GetReportData()  // each solver's history
pool := log.GetSolutionData()  // every pool improvement
fmt.Printf("%d solver(s), %d pool improvement(s)\n", len(report), len(pool))

log.SaveCsv() // MyProblem-<solver>-<id>.csv and MyProblem-pool.csv
```

`Solve` resets the logger at the start of every run, so the results describe the last run. `Solve` never calls `Close`: defer it yourself, and read or save the results before it runs.

## Built-in loggers

| | `logger/stdout` | `logger/channel` |
|---|---|---|
| Output | Plain lines: a solver's progress when its best cost improves, otherwise at most once per print interval | A live dashboard, redrawn in place |
| Suits | Batch jobs, CI, output redirected to a file | Interactive terminal runs |
| History buffer per solver | 512 KB | 128 KB |

| Constructor | Level |
|---|---|
| `stdout.DefaultLogger(problemName)` | `logger.DefaultLogLevel` (`INFO`) |
| `stdout.NewLogger(problemName, level)` | `level` |
| `channel.DefaultLogger(problemName)` | `logger.DefaultLogLevel` (`INFO`) |
| `channel.NewLoggerLevel(problemName, level)` | `level` |
| `channel.NewLogger(problemName, level, bufferSize)` | `level`, with `bufferSize` slots in the update channel (default 15) |

All return the concrete type (`*stdout.Log`, `*channel.Log`), which has the setters below. The [`configuration` package](../configuration/README.md) picks one from the YAML `LogType` key (`PRINT` or `CHANNEL`) and returns it as a `logger.Logger`; type-assert to reach the setters (see [Saving to CSV](#saving-to-csv)).

### `stdout`

- Prints a solver's progress when its best cost improves, otherwise at most every 500 ms. Verbose messages are always printed.
- `SetPrintInterval(timeMilliseconds int)`: minimum time between two prints of the same solver.

### `channel`

- Draws the dashboard in the terminal's alternate screen. When the run ends, it leaves the alternate screen and prints the final frame, then each solver's done message, so the result stays in the scrollback.
- `SetTicker(timeMilliseconds int)`: dashboard refresh interval (default 500 ms); values below 300 ms are ignored.
- `SetNumPoolMessages(num int)`: how many recent pool entries the dashboard shows.
- `SetNumVerboseMessages(num int)`: how many recent verbose messages per solver the dashboard shows.
- `Start`, `Shutdown`, `Print`, `Finish`, `WorkersPrint` and `GetTicker` are called by `Solve`; you don't call them.

## Levels

| Level | Prints |
|---|---|
| `logger.SILENT` | Nothing. |
| `logger.INFO` | Progress (default, `logger.DefaultLogLevel`). |
| `logger.VERBOSE` | Progress and `Verbose` events, such as population restarts. |

| Function | Does |
|---|---|
| `GetLogLevel(name) Level` | Parses `"SILENT"`, `"INFO"` or `"VERBOSE"`, case-insensitive; unknown values give `INFO`. Used for the YAML `LogLevel` key. |
| `GetLevel(level) (Level, error)` | The same, but unknown values give `SILENT` and an error. |
| `GetLevelString(level) string` | `"Silent"`, `"Info"` or `"Verbose"`. |
| `GetLogType(label) LogType` | Parses `"PRINT"` or `"CHANNEL"` into `logger.PRINT` or `logger.CHANNEL`; unknown values give `PRINT`. |

## Results

| Method | Returns |
|---|---|
| `GetReportData() []SolverInformation` | Each solver's name, id and history: one `Data` per `Register` call. |
| `GetSolutionData() []SolutionData` | Every solution added to the pool, in order. |

```go
type SolverInformation struct {
    Name        string
    Id          int
    Performance []Data
}

type Data struct {
    LocalCost int
    BestCost  int
    Time      float64 // seconds since the run started
}

type SolutionData struct {
    Cost int
    Time float64
}
```

## Saving to CSV

`*stdout.Log` and `*channel.Log` have `SaveCsv(filename ...string)`. It writes `<prefix>-<solver>-<id>.csv` (header `best,local,time`) per solver and `<prefix>-pool.csv` (header `cost,time`), where the prefix is the argument, or the problem name without one. It streams each solver's history file into its CSV, so it never loads the whole history in memory. It calls `log.Fatal` on an I/O error.

`SaveCsv` is not part of the `Logger` interface. With the `logger.Logger` returned by `configuration.CreateSolver`, type-assert to reach it:

```go
solver, rkoLog := configuration.CreateSolver("MyProblem", env, cfg)
defer rkoLog.Close()

solver.Solve(context.Background())

if saver, ok := rkoLog.(interface{ SaveCsv(filename ...string) }); ok {
    saver.SaveCsv()
}
```

The package also has plain CSV writers:

| Function | Writes |
|---|---|
| `SaveCsvFile(filename, rows)` | `rows` as CSV; calls `log.Fatal` on an I/O error. |
| `SavePoolCsv(poolData, problemName, filename...)` | Header `cost,time`; default file `<problemName>-pool.csv`. |
| `SaveSolverCSV(solverInfo, problemName, filename...)` | Header `best,local,time`; default file `<problemName>-<Name>-<Id>.csv`. |

## History files

Both loggers keep each solver's history in a memory buffer of 24 bytes per record. When the buffer is full, they append it to that solver's history file in one write and reuse the buffer, so a long run doesn't hold its whole history in memory.

- The file is created on the first flush; a run that never fills the buffer never touches the disk. After its first flush, a solver also keeps an encoding buffer of the same size.
- `GetReportData` and `SaveCsv` read the file, then the records still in memory. `GetReportData` returns the whole history in one slice.
- Files are named `rko-history-<problem>-<solver>-<id>-*.bin`, in the system temp directory (`os.TempDir()`). Characters invalid in file names become `_`, and each name part is cut to 32 bytes.
- Each record is 24 bytes, little-endian: local cost (`int64`), best cost (`int64`), time in seconds (`float64`).
- If a solver's file cannot be created or written, that solver keeps the rest of its history in memory and the logger prints one warning to stderr.

| Method | Does |
|---|---|
| `SetHistoryBuffer(bytes int)` | Sets the buffer size per solver. `0` keeps the whole history in memory. |
| `SetHistoryDir(dir string)` | Sets the directory of the history files. `""` means the system temp directory. |
| `Close() error` | Deletes the last run's history files. |

The setters are methods of `*stdout.Log` and `*channel.Log`, not of the `Logger` interface. They take effect at the next `Solve`, so you can call them after creating the solver:

```go
solver, rkoLog := configuration.CreateSolver("MyProblem", env, cfg)
defer rkoLog.Close()

if l, ok := rkoLog.(interface{ SetHistoryBuffer(bytes int) }); ok {
    l.SetHistoryBuffer(1 << 20) // 1 MB per solver
}

solver.Solve(context.Background())
// read GetReportData or call SaveCsv here, before Close
```

### Close

- Each `Solve` deletes the previous run's files; `Close` deletes the last run's. `Solve` and `configuration.CreateSolver` never call it.
- Call it once you have read or saved the results; `defer` right after creating the logger does that.
- It is safe to call twice. Afterwards the history stays in memory, and spilling resumes at the next `Solve`.
- Call `Close` on a `stdout` logger only after `Solve` returns.

## Writing your own logger

Implement both interfaces and pass your `Logger` to any solver constructor.

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

type SolverLogger interface {
    Register(local int, localBest int, time float64, extra string)
    Verbose(message string, timeStamp float64)
}
```

| Method | Called |
|---|---|
| `AddSolutionPool` | When the pool accepts a new best solution. |
| `WorkerDone` | When a metaheuristic finishes, with its final message. |
| `GetLogger` | Once per metaheuristic, at creation, to get its `SolverLogger`. |
| `GetLogLevel` | By `Solve`, to decide what to print. |
| `GetReportData`, `GetSolutionData` | By you, after `Solve`. |
| `Reset` | By `Solve` at the start of every run; empty the previous run's data. |
| `Close` | By you, after reading the results; release files or connections. Return `nil` if there is nothing to release. |
| `Register` | By a metaheuristic, once per iteration: current cost, its best cost, elapsed seconds and optional extra text. |
| `Verbose` | By a metaheuristic, for occasional events such as a population restart. |

Each metaheuristic runs in its own goroutine, so `WorkerDone` and different `SolverLogger`s can be called at the same time.

## See also

- [Project README](../README.md)
- [`configuration` package](../configuration/README.md)
- [`definition` package](../definition/README.md)
- [`metaheuristc` package](../metaheuristc/README.md)
