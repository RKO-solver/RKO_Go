# RKO-Go

RKO-Go is a Go library for combinatorial optimization with **Random Key Optimization (RKO)**. You describe your problem once, by implementing the [`definition.Environment`](definition/README.md) interface, and the library runs several metaheuristics on it in parallel, sharing their best solutions through a common pool.

Based on [A Random-Key optimizer for combinatorial optimization](https://doi.org/10.1007/s10732-025-09568-z).

## Features

- **Seven metaheuristics:** MultiStart, Simulated Annealing (SA), Iterated Local Search (ILS), Variable Neighbourhood Search (VNS), Genetic Algorithm (GA), Biased Random-Key Genetic Algorithm (BRKGA) and Large Neighbourhood Search (LNS).
- **Problem-agnostic:** every metaheuristic works on any `Environment`.
- **Parallel:** each selected metaheuristic runs in its own goroutine; they exchange solutions through a thread-safe pool.
- **Stopping rules:** a time limit, an iteration limit and your own `context.Context`.
- **Configurable** from Go or from a YAML file, with a terminal dashboard or plain-text logging.

## Install

```sh
go get github.com/RKO-solver/rko-go
```

The root package is named `rko`. Requires Go 1.26 or later.

## Quick start

A solution is a vector of random keys in `[0, 1)`. You write how many keys your problem needs, what a vector costs (lower is better) and how to turn it into an answer. This program solves a five-city travelling salesman problem:

```go
package main

import (
	"context"
	"fmt"
	"time"

	rko "github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger/printer"
)

// TSP is a travelling salesman problem: visit every city once and return to the start.
type TSP struct {
	Dist [][]int
}

// NumKeys returns one key per city.
func (t TSP) NumKeys() int { return len(t.Dist) }

// Cost returns the length of the tour encoded by r.
func (t TSP) Cost(r definition.RandomKey) int {
	tour := r.SortedIndex()
	cost := 0
	for i := range tour {
		cost += t.Dist[tour[i]][tour[(i+1)%len(tour)]]
	}
	return cost
}

// Decode returns the tour: cities in ascending order of their keys.
func (t TSP) Decode(r definition.RandomKey) any { return r.SortedIndex() }

// SwapSearch lets the swap local search exchange any two keys.
func (t TSP) SwapSearch() [][2]int { return [][2]int{{0, t.NumKeys()}} }

// ProblemName returns the name that labels the logger and its output files.
func (t TSP) ProblemName() string { return "tsp" }

func main() {
	env := TSP{Dist: [][]int{
		{0, 2, 9, 10, 7},
		{2, 0, 6, 4, 3},
		{9, 6, 0, 8, 5},
		{10, 4, 8, 0, 6},
		{7, 3, 5, 6, 0},
	}}

	log := printer.Default("tsp")
	defer log.Close()

	mh := []rko.MetaHeuristic{rko.ILS, rko.BRKGA}
	solver := rko.CreateDefaultSolverTimeLimit(mh, 5*time.Second, env, log)

	tour := solver.Solve(context.Background()).([]int)
	cost := solver.GetSolutionPool().BestCost()
	fmt.Println("tour:", tour, "cost:", cost)
}
```

After five seconds of progress lines, it prints an optimal tour (its rotation or direction may vary):

```text
tour: [4 0 1 3 2] cost: 26
```

`Solve` returns the best solution found, already passed through your `Decode`. The [`definition` package](definition/README.md) explains each `Environment` method and shows a non-permutation example.

## Run from a YAML file

The [`config` package](config/README.md) builds the same solver from a file, including each metaheuristic's local search and tuning parameters.

`config.yaml`:

```yaml
LogLevel: INFO
LogType: PRINT
TimeLimitSeconds: 60
MaxIterations: 10000
Metaheuristics:
  - "ILS=Swap,Mirror"
  - "BRKGA"
Parameters:
  BRKGA:
    PopulationSize: 100
```

```go
cfgs, err := config.ReadFile("config.yaml")
if err != nil {
	log.Fatal(err)
}

solver, rkoLog, err := cfgs[0].NewSolver(env)
if err != nil {
	log.Fatal(err)
}
defer rkoLog.Close()

tour := solver.Solve(context.Background()).([]int)
```

See the [`config` package](config/README.md) for every key and its default.

## Stopping rules

`Solve(ctx)` stops at whichever comes first:

- the **time limit** expires,
- every metaheuristic has run **`MaxIterations`** iterations,
- **`ctx`** is cancelled or reaches its deadline.

| Constructor | Time limit | Iteration limit |
|---|---|---|
| `rko.CreateDefaultSolver(mh, env, log)` | 150 s | none |
| `rko.CreateDefaultSolverTimeLimit(mh, timeLimit, env, log)` | `timeLimit` | none |
| `rko.CreateFullSolver(log, env, timeLimit, maxIterations, pool, solvers)` | `timeLimit` | `maxIterations` |
| `cfg.NewSolver(env)` ([`config`](config/README.md)) | `TimeLimitSeconds` in YAML (`cfg.TimeLimit`) | `MaxIterations` in YAML (`cfg.MaxIterations`) |

- A zero or negative value means no limit.
- `MaxIterations` is one number for the whole run, counted by each metaheuristic separately: each one stops after that many iterations of its own outer loop (a generation for GA and BRKGA, a temperature step for SA). `Solve` returns once all of them have stopped.
- With neither limit, the solver runs until `ctx` is cancelled. `rko.CreateFullSolver`, and so `cfg.NewSolver`, prints a warning in that case.
- SA also ends on its own once its temperature reaches `TemperatureGoal`.

To stop the search yourself, for example on a shutdown signal, pass your own context. The configured limits still apply:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

result := solver.Solve(ctx)
```

## Choosing metaheuristics

`rko.CreateDefaultSolver` and `rko.CreateDefaultSolverTimeLimit` take a slice of `rko.MetaHeuristic` constants and use each one's default parameters and local search:

| Constant | Metaheuristic |
|---|---|
| `rko.MULTISTART` | MultiStart |
| `rko.SA` | Simulated Annealing |
| `rko.GA` | Genetic Algorithm |
| `rko.BRKGA` | Biased Random-Key Genetic Algorithm |
| `rko.ILS` | Iterated Local Search |
| `rko.VNS` | Variable Neighbourhood Search |
| `rko.LNS` | Large Neighbourhood Search |

Listing a constant twice runs two independent instances. `GRASP`, `VLNS`, `ALNS` and `IPR` are declared but not implemented; they are skipped with a `not implemented yet` message.

To change parameters or local searches, use a [YAML file](config/README.md), or build each solver yourself and pass them to `rko.CreateFullSolver` (see the [`metaheuristc` package](metaheuristc/README.md)).

## Logging

Every constructor takes a [`logger.Logger`](logger/README.md). Two ship with the library:

| Logger | Constructor | Output |
|---|---|---|
| `logger/printer` | `printer.Default(name)`, `printer.New(w, name, opts)` | Plain lines: a solver's progress when its best cost improves, otherwise at most every 500 ms (`Options.PrintInterval`). Suits batch jobs and CI. |
| `logger/dashboard` | `dashboard.New(name, opts)` | A live dashboard in the terminal's alternate screen; the final frame stays in the scrollback. Suits interactive runs. |

Levels are `logger.LevelSilent`, `logger.LevelInfo` (default) and `logger.LevelVerbose`. `Solve` resets the logger at the start of every run, so `Report` and `Pool` describe the last run. Both loggers can save the run to CSV (`SaveCSV`), and spill each solver's history to [temporary files](logger/README.md#history-files) on long runs. `Solve` never calls the logger's `Close`, which deletes those files: `defer` it after creating the logger. See the [`logger` package](logger/README.md).

## Packages

| Package | Use it to |
|---|---|
| `rko` (root) | Create a `Solver` (`CreateDefaultSolver`, `CreateDefaultSolverTimeLimit`, `CreateFullSolver`) and run it (`Solve`, `GetSolutionPool`, `Print`). |
| [`definition`](definition/README.md) | Implement `Environment` for your problem, or `Solver` for your own metaheuristic. |
| [`config`](config/README.md) | Build a solver from a YAML file or a Go struct. |
| [`logger`](logger/README.md) | Choose or implement a logger, and export results to CSV. |
| [`metaheuristc`](metaheuristc/README.md) | Build metaheuristics, local searches and the solution pool by hand. |
| `random` | Create random generators, seeded (`New`) or randomly seeded (`NewRandom`). |

## Maintainers

Lucas Mendes - [GitHub](https://github.com/lucasmends)

## License

Licensed under the Academic Free License v3.0. See [LICENSE](LICENSE).
