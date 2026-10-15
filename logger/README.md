# `logger` Package

The loggers RKO-Go metaheuristics report to. A logger shows progress while a run is going, records each metaheuristic's history and every solution-pool improvement, and saves them to CSV. Two ship with the library, `printer` and `dashboard`; implement `Logger` to write your own.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"time"

	rko "github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/problem"
)

func main() {
	var env problem.Environment // your problem

	log := printer.New(nil, "MyProblem", &printer.Options{Level: logger.LevelVerbose})
	defer log.Close() // deletes the history files, so read the results first

	solver, err := rko.New(env, log,
		rko.WithMetaheuristics(metaheuristic.GA, metaheuristic.SA),
		rko.WithTimeLimit(10*time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}
	result := solver.Solve(context.Background())
	fmt.Printf("Best solution: %+v\n", result)

	fmt.Printf("%d metaheuristic(s), %d pool improvement(s)\n", len(log.Report()), len(log.Pool()))

	if err := log.SaveCSV(""); err != nil { // MyProblem-<metaheuristic>-<id>.csv and MyProblem-pool.csv
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
| Default history buffer per metaheuristic | 512 KiB | 128 KiB |

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

Prints a metaheuristic's progress line on its first iteration, then skips `PrintInterval` iterations after each printed line. Events and pool improvements are printed immediately (at `LevelVerbose`). 

> **Warning:** the `printer` does not lock. Metaheuristic goroutines write to `w` directly, so `w` must be safe for concurrent use. `os.Stdout` and `os.Stderr` are. A `bytes.Buffer`, a `bufio.Writer` and most custom writers are not: wrap them with a mutex.
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
> Nothing else in the `printer` is synchronized either: `MetaheuristicLogger`, `Reset`, `Report`, `Pool`, the CSV and policy methods and `Close` are for before the run (setup) or after it (results). Only the `dashboard` is safe to read during a run.

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

Draws in the terminal's alternate screen and refreshes every `Refresh`. When the run ends it leaves the alternate screen and prints the final frame, then each metaheuristic's `done` line, so the result stays in the scrollback. It writes nothing at `LevelSilent`.

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

- Each metaheuristic gets a block of three lines: iteration, local and best cost, and time. A metaheuristic appears after its first `Register`.
- At `LevelVerbose`, the last `PoolLines` pool improvements are listed at the top, and each block adds one line per event key, holding the latest value of that key.
- `(id) name started` lines are printed before the first frame, so they stay in the scrollback.
- Lines are cut to fit (300 bytes for pool lines, 100 for metaheuristic lines).
- Policies are not shown; print them after the run with `PrintPolicies`, which writes them to `os.Stdout` (see [Policies](#policies)).

#### Running the dashboard

The dashboard implements `logger.Runner`: `Start()` launches the goroutine that collects the metaheuristics' updates and the one that redraws the screen, and returns `stop`. `rko.Solver.Solve` calls `Start` before the metaheuristics and `stop` after, so you do nothing special.

Outside `Solve`, for example in your own loop, call both yourself:

```go
log := dashboard.New("MyProblem", nil)
defer log.Close()

stop := log.Start()
// ... run metaheuristics that use log.MetaheuristicLogger(...)
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
| `PrintInterval` | 5 | 0 | Iterations skipped after each printed progress line; `0` selects the default, negative values print every iteration. |
| `HistoryBuffer` | 512 KiB | | Bytes of history kept in memory per metaheuristic before it spills to a file. Negative keeps all history in memory. |
| `HistoryDir` | system temp directory | | Where history files are created. |

### `dashboard.Options`

| Field | Default | Minimum | Meaning |
|---|---|---|---|
| `Level` | `LevelInfo` | | What to show. |
| `Refresh` | 500 ms | 300 ms | Redraw interval. |
| `BufferSize` | 1024 | 1 | Slots in the update queue between the metaheuristics and the dashboard. |
| `PoolLines` | 3 | 1 | Pool improvements shown (Verbose). |
| `HistoryBuffer` | 128 KiB | | Bytes of history kept in memory per metaheuristic before it spills to a file. Negative keeps all history in memory. |
| `HistoryDir` | system temp directory | | Where history files are created. |

## Reading results

Both loggers implement `logger.Recorder`, which is what `config.NewSolver` returns:

| Method | From | Returns |
|---|---|---|
| `Report() []logger.MetaheuristicReport` | `logger.Reporter` | Each metaheuristic's name, id and history. |
| `Pool() []logger.PoolEntry` | `logger.Reporter` | Every solution added to the pool, in order. |
| `SaveCSV(prefix) error` | `logger.CSVSaver` | One CSV file per metaheuristic and one for the pool. |
| `PrintPolicies() error` | `logger.PolicySaver` | Writes each recorded policy as text to the logger's output (see [Policies](#policies)). |
| `SavePoliciesCSV(prefix) error` | `logger.PolicySaver` | One CSV file per metaheuristic that recorded a policy. |
| `Close() error` | `io.Closer` | Deletes the history files (see [History files](#history-files-and-close)). |

```go
type MetaheuristicReport struct {
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

Read the results after `Solve` returns. Only the `dashboard` makes `Report`, `Pool`, `SaveCSV` and the policy methods safe during a run; the `printer` does not. For any other format, write it from `Report` and `Pool`.

### CSV

| CSV | Header | Columns |
|---|---|---|
| Metaheuristic | `best,local,time` | Best cost, local cost, seconds since the run started (8 decimals). One row per `Register` call. |
| Pool | `cost,time` | Cost and seconds, one row per pool improvement. |

`SaveCSV(prefix)` writes `<prefix>-<metaheuristic>-<id>.csv` for every metaheuristic and `<prefix>-pool.csv`:

- `prefix` is used as is, so it may include a directory (`out/run1`), which must exist.
- An empty `prefix` uses the problem name, with path separators and characters that file systems reject replaced by `_` and cut to 32 bytes; `rko` if that leaves nothing.
- The metaheuristic part of the name is cleaned the same way.
- It tries every file and returns all errors joined.

### Policies

A metaheuristic that tunes its parameters with a [`tuning`](../metaheuristic/tuning/README.md) tuner can record the policy the tuner learned, by calling `MetaheuristicLogger.Policy` once at the end of its run. A policy is a `logger.Table`:

```go
type Table struct {
	Header []string
	Rows   [][]string
}
```

The logger never prints policies by itself; call these methods after the run, as with `SaveCSV`:

```go
if err := log.PrintPolicies(); err != nil {
	fmt.Println(err)
}
if err := log.SavePoliciesCSV(""); err != nil { // MyProblem-<metaheuristic>-<id>-policy.csv
	fmt.Println(err)
}
```

`PrintPolicies()` writes to the logger's output: the `w` given to `printer.New` (`os.Stdout` if `nil`), or `os.Stdout` for `dashboard`. For each metaheuristic that recorded a policy and in id order, it writes a `(id) name policy` heading and then the table in aligned columns, with a blank line between metaheuristics. `printer` and `dashboard` use the same format. Here GA tunes `PopulationSize: [100, 200]` and `MutationAlpha: [0.005, 0.01]`, and BRKGA `PopulationSize: [100, 200]` and `EliteRatio: [0.1, 0.2]`:

```text
(0) GA policy
state  p0   p1    p2     p3   move       next  q
0      100  0.95  0.005  100  stay       0     -0.05672288500153648
1      200  0.95  0.005  100  p0 -> 100  0     -0.07584516941366984
2      100  0.95  0.01   100  stay       2     -0.03486735460061637
3      200  0.95  0.01   100  p0 -> 100  2     -0.06383021884968096

(1) BRKGA policy
state  p0   p1   p2    p3    p4     p5   move       next  q
0      100  0.1  0.05  0.95  0.005  100  p0 -> 200  1     0.14475877694232625
1      200  0.1  0.05  0.95  0.005  100  p0 -> 100  0     0.04066887036137811
2      100  0.2  0.05  0.95  0.005  100  p0 -> 200  3     0.040684180716951374
3      200  0.2  0.05  0.95  0.005  100  p0 -> 100  2     0.040708076771112024
```

The columns depend on the tuner; for `QLearning`, see [Reading the learned policy](../metaheuristic/tuning/README.md#reading-the-learned-policy). The built-in metaheuristics give one `p` column per tunable `Config` field, even one with a single value, in the order listed in [Learned policy](../metaheuristic/README.md#learned-policy).

`SavePoliciesCSV(prefix)` writes `<prefix>-<metaheuristic>-<id>-policy.csv` for each metaheuristic that recorded a policy: the header row, then one row per table row. `prefix` and the metaheuristic part are handled as in `SaveCSV`, and an empty `prefix` uses the problem name.

- If no metaheuristic recorded a policy, both methods write nothing and return `nil`.
- `Reset` (so every `Solve`) clears the policies along with the history.
- The `dashboard` stores a policy directly, not through its update queue, so `Policy` never blocks and works without `Start`.

## Writing a metaheuristic

A metaheuristic gets a `logger.MetaheuristicLogger` from `Logger.MetaheuristicLogger(name)` and reports to it:

```go
type MetaheuristicLogger interface {
	Register(local, best, iteration int, elapsed time.Duration)
	Event(a slog.Attr, elapsed time.Duration)
	Policy(t Table)
}
```

| Method | Call it | Notes |
|---|---|---|
| `Register` | Once per iteration, with the current cost, the best cost, the iteration number and the time since the run started. | Takes no strings and allocates nothing. Recorded at every level. |
| `Event` | For occasional events, such as a restart or a temperature change. | Shown only at `LevelVerbose`. |
| `Policy` | Once at the end of the run, if the metaheuristic tunes its parameters, with the tuner's learned policy. | Kept without copying, so don't modify `t` afterwards. Never shown during the run (see [Policies](#policies)). |

```go
ctx, counter := metaheuristic.NewSolverContext(ctx)
defer counter.Stop()

for ; ctx.Err() == nil; counter.Step() {
	// ... improve localSolution and bestSolution ...
	elapsed := time.Since(start)

	log.Register(localSolution.Cost, bestSolution.Cost, counter.Iteration(), elapsed)
	log.Event(slog.Float64("temperature", temperature), elapsed)
}
```

Build the attribute with a typed constructor (`slog.Int`, `slog.Float64`, `slog.String`, `slog.Bool`, `slog.Duration`) and a constant key. That allocates nothing, and the value is formatted only when shown. The `dashboard` keeps the latest value of each key, so use a stable key per quantity (`"temperature"`, not `"temperature 10"`).

`slog.GroupValue` joins several values into one event, but allocates; use it only for rare events. The built-in metaheuristics use it for `stalled`, `improved` and SA's `restart`.

The built-in metaheuristics emit these events:

| Metaheuristic | Key | Value | Fires |
|---|---|---|---|
| GA, BRKGA | `stalled` | group: `generation`, `no-improvement` (ints) | After each generation that did not improve the best solution. |
| GA, BRKGA | `restart` | string | When the population restarts. |
| SA | `temperature` | float | After every temperature step. |
| SA | `reheat` | int | At each reheat, with the reheat number (`1` for the first). |
| SA | `restart` | group: `iteration` (int) | At each restart. The best cost passed to `Register` then restarts too, so it can go up. |
| VNS | `improved` | group: `iteration`, `local` (ints) | When a neighbour beats the best solution. |

MultiStart, ILS and LNS emit none.

A `MetaheuristicLogger` is used by one goroutine at a time, so use one per metaheuristic. With the `dashboard`, `Register` and `Event` block when the update queue is full.

The built-in metaheuristics also report pool improvements through `solution.Pool.Add`, which calls `Logger.AddSolutionPool`. See the [`problem`](../problem/README.md) and [`metaheuristic`](../metaheuristic/README.md) packages.

## History files and Close

Each metaheuristic keeps its history in a memory buffer of 24 bytes per record. When the buffer is full, the logger appends it to that metaheuristic's file in one write and reuses the buffer, so a long run does not hold its whole history in memory.

- The file is created on the first spill; a run that never fills the buffer never touches the disk.
- `Report` and `SaveCSV` read the file, then the records still in memory.
- Files are named `rko-history-<problem>-<metaheuristic>-<id>-*.bin`, in `HistoryDir` or the system temp directory. Each name part is cleaned of characters invalid in file names and cut to 32 bytes.
- If a file cannot be created or written, that metaheuristic keeps the rest of its history in memory and one warning is printed to stderr.
- `HistoryBuffer` negative keeps everything in memory and creates no files.

`Close` deletes these files. Each `Solve` also deletes the previous run's files when it resets the logger.

- Read or save the results first: once `Close` has run, the spilled part of the history is gone. `defer log.Close()` right after creating the logger does that.
- It is safe to call twice. Until the next `Solve` (or `Reset`), the logger keeps new history in memory.
- Call it after `Solve` returns.

## Writing your own `Logger`

Implement `logger.Logger` and pass it to any solver constructor. It holds only what the solver calls:

```go
type Logger interface {
	MetaheuristicLogger(name string) MetaheuristicLogger
	AddSolutionPool(cost int, elapsed time.Duration)
	WorkerStarted(id int, name string)
	WorkerDone(r WorkerResult)
	Reset()
}
```

| Method | Called |
|---|---|
| `MetaheuristicLogger` | By each metaheuristic when it is created, to get its `MetaheuristicLogger`. Number metaheuristics from 0 in call order: that is `MetaheuristicReport.ID`. |
| `AddSolutionPool` | When the pool accepts a new best solution. |
| `WorkerStarted`, `WorkerDone` | By `Solve`, when a metaheuristic starts and when it finishes (with its id, name, final cost and time). |
| `Reset` | By `Solve` at the start of every run; empty the previous run's data, keeping the registered metaheuristics. |

The `MetaheuristicLogger` you return must implement all three of its methods, including `Policy` (see [Writing a metaheuristic](#writing-a-metaheuristic)); if your logger doesn't keep policies, `Policy` can do nothing.

Everything else is optional; implement only what your logger offers:

| Interface | Methods | For |
|---|---|---|
| `logger.Reporter` | `Report`, `Pool` | Returning the recorded run. |
| `logger.CSVSaver` | `SaveCSV` | Saving the run as CSV files. |
| `logger.PolicySaver` | `PrintPolicies`, `SavePoliciesCSV` | Printing and saving the recorded policies. |
| `io.Closer` | `Close` | Releasing files or other resources. |
| `logger.Recorder` | all of the above plus `Logger` | What `printer` and `dashboard` implement. |

Metaheuristics run in their own goroutines, so `WorkerStarted`, `WorkerDone`, `AddSolutionPool` and different `MetaheuristicLogger`s can be called at the same time.

If your logger needs a loop during the run, also implement `logger.Runner`:

```go
type Runner interface{ Start() (stop func()) }
```

`Solve` calls `Start` after `Reset` and before it starts the metaheuristics, and the returned `stop` after every metaheuristic has finished.

## See also

- [Project README](../README.md)
- [`config` package](../config/README.md)
- [`problem` package](../problem/README.md)
- [`metaheuristic` package](../metaheuristic/README.md)
