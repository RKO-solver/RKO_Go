# `logger` Package

The loggers RKO-Go solvers report to. A logger shows progress while a run is going, records each solver's history and every solution-pool improvement, and saves them to CSV. Two ship with the library, `printer` and `dashboard`; implement `Logger` to write your own.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"time"

	rko "github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/printer"
)

func main() {
	var env definition.Environment // your problem

	log := printer.New(nil, "MyProblem", &printer.Options{Level: logger.LevelVerbose})
	defer log.Close() // deletes the history files, so read the results first

	solver := rko.New(env, log, []rko.Metaheuristic{rko.GA, rko.SA}, rko.WithTimeLimit(10*time.Second))
	result := solver.Solve(context.Background())
	fmt.Printf("Best solution: %+v\n", result)

	fmt.Printf("%d solver(s), %d pool improvement(s)\n", len(log.Report()), len(log.Pool()))

	if err := log.SaveCSV(""); err != nil { // MyProblem-<solver>-<id>.csv and MyProblem-pool.csv
		fmt.Println(err)
	}
}
```

`Solve` resets the logger at the start of every run, so the results describe the last run. `Solve` never calls `Close`: `defer` it yourself, and read or save the results before it runs.

## Printer or dashboard

| | `logger/printer` | `logger/dashboard` |
|---|---|---|
| Output | One plain line per progress print | A live screen, redrawn in place |
| Suits | Batch jobs, CI, output redirected to a file or any `io.Writer` | Interactive terminal runs |
| Writes to | Any `io.Writer` (`os.Stdout` if `nil`) | `os.Stdout` only |
| Default history buffer per solver | 512 KiB | 128 KiB |

| Constructor | Does |
|---|---|
| `printer.New(w, problem, opts)` | Writes to `w`, which must be safe for concurrent use; a `nil` `w` means `os.Stdout`. |
| `printer.Default(problem)` | `printer.New(os.Stdout, problem, nil)`. |
| `dashboard.New(problem, opts)` | Draws to the terminal. |

`problem` labels the output and names the default CSV files. `opts` may be `nil`. Both return a concrete `*Log` that implements `logger.Recorder`.

Select one by name with `logger.Type`, as the [`config` package](../config/README.md) does with the YAML `LogType` key:

| Constant | Text | Logger |
|---|---|---|
| `logger.TypePrint` (zero value) | `PRINT` | `printer` |
| `logger.TypeDashboard` | `DASHBOARD` | `dashboard` |

`logger.ParseType` ignores case. `Type` also implements `String`.

### `printer`

Prints a solver's progress line when its best cost improves, otherwise at most once per `PrintInterval`. Events and pool improvements are printed immediately (at `LevelVerbose`). 

> **Warning:** the `printer` does not lock. Solver goroutines write to `w` directly, so `w` must be safe for concurrent use. `os.Stdout` and `os.Stderr` are. A `bytes.Buffer`, a `bufio.Writer` and most custom writers are not: wrap them with a mutex.
>
> ```go
> type lockedWriter struct {
> 	mu sync.Mutex
> 	w  io.Writer
> }
>
> func (l *lockedWriter) Write(p []byte) (int, error) {
> 	l.mu.Lock()
> 	defer l.mu.Unlock()
> 	return l.w.Write(p)
> }
>
> log := printer.New(&lockedWriter{w: &buf}, "MyProblem", nil)
> ```
>
> Nothing else in the `printer` is synchronized either: `SolverLogger`, `Reset`, `Report`, `Pool`, the CSV methods and `Close` are for before the run (setup) or after it (results). Only the `dashboard` is safe to read during a run.

```text
(0) GA started
(0) GA iteration=152 local=10 best=9 time=1.234s
(0) GA stalled=[generation=153 no-improvement=3] time=1.250s
pool cost=9 time=1.234s
(0) GA done cost=9 time=12.345s
```

| Line | Level |
|---|---|
| `(id) name started`, `(id) name done cost=... time=...` | Info and Verbose |
| `(id) name iteration=... local=... best=... time=...` | Info and Verbose |
| `(id) name key=value time=...` (an event) | Verbose |
| `pool cost=... time=...` (a pool improvement) | Verbose |

### `dashboard`

Draws in the terminal's alternate screen and refreshes every `Refresh`. When the run ends it leaves the alternate screen and prints the final frame, then each solver's `done` line, so the result stays in the scrollback. It writes nothing at `LevelSilent`.

```text

pool cost=9 time=1.234s
pool cost=8 time=2.101s

(0) GA: iteration 152
    Local 10 Best 9
    Time 1.234s
    stalled=[generation=153 no-improvement=3] time=1.250s

(1) SA: iteration 4210
    Local 12 Best 8
    Time 2.101s
```

- Each solver gets a block of three lines: iteration, local and best cost, and time. A solver appears after its first `Register`.
- At `LevelVerbose`, the last `PoolLines` pool improvements are listed at the top, and each block adds one line per event key, holding the latest value of that key.
- `(id) name started` lines are printed before the first frame, so they stay in the scrollback.
- Lines are cut to fit (300 bytes for pool lines, 100 for solver lines).

#### Running the dashboard

The dashboard implements `logger.Runner`: `Start()` launches the goroutine that collects the solvers' updates and the one that redraws the screen, and returns `stop`. `rko.Solver.Solve` calls `Start` before the solvers and `stop` after, so you do nothing special.

Outside `Solve`, for example in your own loop, call both yourself:

```go
log := dashboard.New("MyProblem", nil)
defer log.Close()

stop := log.Start()
// ... run solvers that use log.SolverLogger(...)
stop() // returns after the final output is written
```

Call `stop` once per `Start`, and don't call `Start` or `Reset` while the dashboard is running. Without `Start`, `Register` blocks once `BufferSize` messages are queued.

## Levels

| Constant | Text | Shows |
|---|---|---|
| `logger.LevelSilent` (`-1`) | `SILENT` | Nothing. |
| `logger.LevelInfo` (`0`, zero value, default) | `INFO` | Progress, and worker start and done. |
| `logger.LevelVerbose` (`1`) | `VERBOSE` | Also events and pool improvements. |

The history is recorded at every level, so `Report`, `Pool` and the CSV files are complete even at `LevelSilent`.

`logger.DefaultLevel` is `LevelInfo`. `logger.ParseLevel` ignores case and returns an error for an unknown name. `Level` also implements `String`.

## Options

Pass `nil` or a zero `Options` for the defaults. A zero field means the default; a value below the minimum is raised to it.

### `printer.Options`

| Field | Default | Minimum | Meaning |
|---|---|---|---|
| `Level` | `LevelInfo` | | What to print. |
| `PrintInterval` | 500 ms | 10 ms | Minimum time between two prints of the same solver, unless its best cost improved. |
| `HistoryBuffer` | 512 KiB | | Bytes of history kept in memory per solver before it spills to a file. Negative keeps all history in memory. |
| `HistoryDir` | system temp directory | | Where history files are created. |

### `dashboard.Options`

| Field | Default | Minimum | Meaning |
|---|---|---|---|
| `Level` | `LevelInfo` | | What to show. |
| `Refresh` | 500 ms | 300 ms | Redraw interval. |
| `BufferSize` | 1024 | 1 | Slots in the update queue between the solvers and the dashboard. |
| `PoolLines` | 3 | 1 | Pool improvements shown (Verbose). |
| `HistoryBuffer` | 128 KiB | | Bytes of history kept in memory per solver before it spills to a file. Negative keeps all history in memory. |
| `HistoryDir` | system temp directory | | Where history files are created. |

## Reading results

Both loggers implement `logger.Recorder`, which is what `config.NewSolver` returns:

| Method | From | Returns |
|---|---|---|
| `Report() []logger.SolverReport` | `logger.Reporter` | Each solver's name, id and history. |
| `Pool() []logger.PoolEntry` | `logger.Reporter` | Every solution added to the pool, in order. |
| `SaveCSV(prefix) error` | `logger.CSVSaver` | One CSV file per solver and one for the pool. |
| `Close() error` | `io.Closer` | Deletes the history files (see [History files](#history-files-and-close)). |

```go
type SolverReport struct {
	Name    string
	ID      int
	Samples []Sample
}

type Sample struct {
	Local, Best int
	Elapsed     time.Duration // since the run started
}

type PoolEntry struct {
	Cost    int
	Elapsed time.Duration
}
```

Read the results after `Solve` returns. Only the `dashboard` makes `Report`, `Pool` and `SaveCSV` safe during a run; the `printer` does not. For any other format, write it from `Report` and `Pool`.

### CSV

| CSV | Header | Columns |
|---|---|---|
| Solver | `best,local,time` | Best cost, local cost, seconds since the run started (8 decimals). One row per `Register` call. |
| Pool | `cost,time` | Cost and seconds, one row per pool improvement. |

`SaveCSV(prefix)` writes `<prefix>-<solver>-<id>.csv` for every solver and `<prefix>-pool.csv`:

- `prefix` is used as is, so it may include a directory (`out/run1`), which must exist.
- An empty `prefix` uses the problem name, with path separators and characters that file systems reject replaced by `_` and cut to 32 bytes; `rko` if that leaves nothing.
- The solver part of the name is cleaned the same way.
- It tries every file and returns all errors joined.

## Writing a metaheuristic

A metaheuristic gets a `logger.SolverLogger` from `Logger.SolverLogger(name)` and reports to it:

```go
type SolverLogger interface {
	Register(local, best, iteration int, elapsed time.Duration)
	Event(a slog.Attr, elapsed time.Duration)
}
```

| Method | Call it | Notes |
|---|---|---|
| `Register` | Once per iteration, with the current cost, the best cost, the iteration number and the time since the run started. | Takes no strings and allocates nothing. Recorded at every level. |
| `Event` | For occasional events, such as a restart or a temperature change. | Shown only at `LevelVerbose`. |

```go
ctx, counter := definition.NewSolverContext(ctx)
defer counter.Stop()

for ; ctx.Err() == nil; counter.Step() {
	// ... improve localSolution and bestSolution ...
	elapsed := time.Since(start)

	log.Register(localSolution.Cost, bestSolution.Cost, counter.Iteration(), elapsed)
	log.Event(slog.Float64("temperature", temperature), elapsed)
}
```

Build the attribute with a typed constructor (`slog.Int`, `slog.Float64`, `slog.String`, `slog.Bool`, `slog.Duration`) and a constant key. That allocates nothing, and the value is formatted only when shown. The `dashboard` keeps the latest value of each key, so use a stable key per quantity (`"temperature"`, not `"temperature 10"`).

`slog.GroupValue` joins several values into one event, but allocates; use it only for rare events. The built-in metaheuristics use it for `stalled` and `improved`.

The built-in metaheuristics emit these events:

| Metaheuristic | Key | Value | Fires |
|---|---|---|---|
| GA, BRKGA | `stalled` | group: `generation`, `no-improvement` (ints) | After each generation that did not improve the best solution. |
| GA, BRKGA | `restart` | string | When the population restarts. |
| SA | `temperature` | float | After every temperature step. |
| SA | `reheat` | int | At each reheat, with the reheat number (`1` for the first). |
| VNS | `improved` | group: `iteration`, `local` (ints) | When a neighbour beats the best solution. |

MultiStart, ILS and LNS emit none.

A `SolverLogger` is used by one goroutine at a time, so use one per solver. With the `dashboard`, `Register` and `Event` block when the update queue is full.

The built-in solvers also report pool improvements through `solution.Pool.Add`, which calls `Logger.AddSolutionPool`. See the [`definition`](../definition/README.md) and [`metaheuristic`](../metaheuristic/README.md) packages.

## History files and Close

Each solver keeps its history in a memory buffer of 24 bytes per record. When the buffer is full, the logger appends it to that solver's file in one write and reuses the buffer, so a long run does not hold its whole history in memory.

- The file is created on the first spill; a run that never fills the buffer never touches the disk.
- `Report` and the CSV methods read the file, then the records still in memory.
- Files are named `rko-history-<problem>-<solver>-<id>-*.bin`, in `HistoryDir` or the system temp directory. Each name part is cleaned of characters invalid in file names and cut to 32 bytes.
- If a file cannot be created or written, that solver keeps the rest of its history in memory and one warning is printed to stderr.
- `HistoryBuffer` negative keeps everything in memory and creates no files.

`Close` deletes these files. Each `Solve` also deletes the previous run's files when it resets the logger.

- Read or save the results first: once `Close` has run, the spilled part of the history is gone. `defer log.Close()` right after creating the logger does that.
- It is safe to call twice. Until the next `Solve` (or `Reset`), the logger keeps new history in memory.
- Call it after `Solve` returns.

## Writing your own `Logger`

Implement `logger.Logger` and pass it to any solver constructor. It holds only what the solver calls:

```go
type Logger interface {
	SolverLogger(name string) SolverLogger
	AddSolutionPool(cost int, elapsed time.Duration)
	WorkerStarted(id int, name string)
	WorkerDone(r WorkerResult)
	Reset()
}
```

| Method | Called |
|---|---|
| `SolverLogger` | By each metaheuristic when it is created, to get its `SolverLogger`. Number solvers from 0 in call order: that is `SolverReport.ID`. |
| `AddSolutionPool` | When the pool accepts a new best solution. |
| `WorkerStarted`, `WorkerDone` | By `Solve`, when a metaheuristic starts and when it finishes (with its id, name, final cost and time). |
| `Reset` | By `Solve` at the start of every run; empty the previous run's data, keeping the registered solvers. |

Everything else is optional; implement only what your logger offers:

| Interface | Methods | For |
|---|---|---|
| `logger.Reporter` | `Report`, `Pool` | Returning the recorded run. |
| `logger.CSVSaver` | `SaveCSV` | Saving the run as CSV files. |
| `io.Closer` | `Close` | Releasing files or other resources. |
| `logger.Recorder` | all of the above plus `Logger` | What `printer` and `dashboard` implement. |

Metaheuristics run in their own goroutines, so `WorkerStarted`, `WorkerDone`, `AddSolutionPool` and different `SolverLogger`s can be called at the same time.

If your logger needs a loop during the run, also implement `logger.Runner`:

```go
type Runner interface{ Start() (stop func()) }
```

`Solve` calls `Start` after `Reset` and before it starts the solvers, and the returned `stop` after every solver has finished.

## See also

- [Project README](../README.md)
- [`config` package](../config/README.md)
- [`definition` package](../definition/README.md)
- [`metaheuristic` package](../metaheuristic/README.md)
