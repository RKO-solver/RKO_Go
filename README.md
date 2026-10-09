# RKO-Go

RKO-Go is a Go library for solving optimization problems with Random Key Optimization (RKO) and a set of metaheuristics. It is problem-agnostic: implement the [`definition.Environment`](definition/README.md) interface for your problem, and every metaheuristic in the library can solve it.

> **Reference:**
> This library is based on the concepts and methods described in the article:
> [Random Key Optimization: A Unified Framework for Metaheuristics](https://doi.org/10.48550/arXiv.2411.04293)

## Features

- **Metaheuristics implemented:**
  - MultiStart
  - Simulated Annealing (SA)
  - Iterated Local Search (ILS)
  - Variable Neighborhood Search (VNS)
  - Genetic Algorithm (GA)
  - Biased Random Key Genetic Algorithm (BRKGA)
  - Large Neighborhood Search (LNS)

- **Problem-agnostic:** implement one interface and every metaheuristic works on your problem.
- **Parallel by default:** the solver runs the selected metaheuristics concurrently, sharing solutions through a thread-safe pool.
- **Context-driven:** a `context.Context` controls when the search stops.
- **Configurable from YAML or Go**, with pluggable logging.

## Getting Started

### 1. Install

```sh
go get github.com/RKO-solver/rko-go
```

### 2. Implement the Environment interface

```go
type Environment interface {
    NumKeys() int
    Cost(r RandomKey) int
    Decode(r RandomKey) any
    SwapSearch() [][2]int
}
```

- `NumKeys()` — number of keys (variables) in your problem.
- `Cost(r)` — objective value of a solution. **The library minimizes this value.**
- `Decode(r)` — converts random keys into your problem's representation.
- `SwapSearch()` — `(start, end)` index pairs the swap local search may exchange.

See the [`definition` package](definition/README.md) for a complete, compiling example.

### 3. Create and run a solver

```go
import (
    "context"
    "fmt"

    "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/logger/channel"
)

env := MyEnvironment{}                        // your implementation
mh := []rko.MetaHeuristic{rko.GA, rko.SA}     // metaheuristics to run in parallel
log := channel.DefaultLogger("MyProblem")

solver := rko.CreateDefaultSolver(mh, env, log)

result := solver.Solve(context.Background())
fmt.Printf("Best solution: %+v\n", result)
```

`Solve` returns the best solution found, already decoded by `Environment.Decode`.

## Time limits and cancellation

The search stops when its time limit expires, when its iteration limit is reached, or when the context is cancelled — whichever comes first.

An iteration is one outer-loop pass of a metaheuristic: one generation for GA and BRKGA, one temperature step for SA. `Solve` applies the time limit; every limit is optional.

| Entry point | Time limit | Iteration limit |
|---|---|---|
| `rko.CreateDefaultSolver(mh, env, log)` | 150 seconds | none |
| `rko.CreateDefaultSolverTimeLimit(mh, limit, env, log)` | the limit you pass | none |
| `rko.CreateFullSolver(log, env, timeLimit, maxIterations, pool, solvers)` | `timeLimit` | `maxIterations` |
| `configuration.CreateSolver(...)` | the `TimeLimitSeconds` key in the solver YAML | the `MaxIterations` key in the solver YAML |

A non-positive limit means not set. With neither limit, the solver runs until its context is cancelled.

To choose the limit explicitly:

```go
solver := rko.CreateDefaultSolverTimeLimit(mh, 30*time.Second, env, log)
result := solver.Solve(context.Background())
```

To stop the search yourself — for example on a shutdown signal — pass your own context to `Solve`. Both limits still apply:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result := solver.Solve(ctx)
```

## Selecting metaheuristics

`rko.CreateDefaultSolver` takes a slice of `rko.MetaHeuristic` constants: `MULTISTART`, `SA`, `GA`, `VNS`, `ILS`, `BRKGA`, `LNS`. Listing the same one twice runs two independent instances of it in parallel.

For per-metaheuristic parameters and local search selection, use the [`configuration` package](configuration/README.md) instead, which builds the same solver from YAML files.

## Logging

Track progress by passing any implementation of the [`logger.Logger`](logger/README.md) interface:

```go
type Logger interface {
    AddSolutionPool(cost int, time float64)
    WorkerDone(message string)
    GetLogger(name string) SolverLogger
    GetLogLevel() Level
    GetReportData() []SolverInformation
    GetSolutionData() []SolutionData
    Reset()
}
```

Two loggers ship with the library:

- [`logger/channel`](logger/README.md) — live-refreshing terminal dashboard, for interactive runs.
- [`logger/stdout`](logger/README.md) — plain lines as they happen, for batch and CI runs.

Both can export the run to CSV. See the [`logger` package](logger/README.md).

## Documentation

| Package | Contents |
|---|---|
| [`definition`](definition/README.md) | The `Environment` interface you implement, plus `Solver`, `Result` and `RandomKey`. |
| [`configuration`](configuration/README.md) | YAML and Go configuration, and the solver factories. |
| [`logger`](logger/README.md) | Logging interfaces, the two built-in loggers, and CSV export. |
| [`metaheuristc`](metaheuristc/README.md) | The metaheuristic implementations, their parameters, and local searches. |

Root package entry points: `CreateDefaultSolver`, `CreateDefaultSolverTimeLimit` and `CreateFullSolver` in [`create.go`](create.go); `Solve` and `GetSolutionPool` in [`definition.go`](definition.go).

## Maintainers

Lucas Mendes - [GitHub](https://github.com/lucasmends)

## License

This project is licensed under the Academic Free License v3.0. See [LICENSE](LICENSE)
