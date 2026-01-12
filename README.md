# RKO-Go

RKO-Go is a Go library for solving optimization problems using Random Key Optimization (RKO) and several metaheuristics. The library is designed to be problem-agnostic: you only need to implement the [`definition.Environment`](definition/definition.go) interface for your specific problem, and you can leverage a variety of metaheuristics already implemented in this package.

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

- **Problem-agnostic:** Just implement the [`definition.Environment`](definition/definition.go) interface for your problem.
- **Thread-safe solution pool** and flexible logging.
- **Easy to extend** with new metaheuristics or local search strategies.
- **Parallel metaheuristics:** Use [`CreateDefaultSolver`](create.go) to run multiple metaheuristics in parallel, sharing solutions between them.

## Logging

You can track the progress of your optimization by providing a custom logger that implements the [`logger.Logger`](logger/definition.go) interface:

```go
type Logger interface {
    AddSolutionPool(cost int, time float64)
    WorkerDone(message string)
    GetLogger(name string) SolverLogger
    GetLogLevel() Level
    GetReportData() []SolverInformation
    GetSolutionData() []SolutionData
}

type SolverLogger interface {
    Register(local int, localBest int, time float64, extra string)
    Verbose(message string, timeStamp float64)
}
```

### Built-in Loggers

- [`logger/channel`](logger/channel/create.go): Provides detailed logging with communication channels for real-time progress tracking.
- [`logger/stdout`](logger/stdout/create.go): Prints progress and messages directly to the standard output.

You can implement your own logger by satisfying the [`logger.Logger`](logger/definition.go) interface.

## Getting Started

### 1. Install

```sh
go get github.com/RKO-solver/rko-go
```

### 2. Implement the Environment Interface

To use RKO-Go, you must implement the [`definition.Environment`](definition/definition.go) interface:

```go
type Environment interface {
    NumKeys() int
    Cost(r RandomKey) int
    Decode(r RandomKey) any
    SwapSearch() [][2]int
}
```

- `NumKeys()` returns the number of keys (variables) in your problem.
- `Cost(r RandomKey)` evaluates the cost (objective function) for a solution.
- `Decode(r RandomKey)` converts a solution from random keys to your problem's representation.
- `SwapSearch()` returns the pairs (start, end) indices for performing multiple swaps during local search.

### 3. Create and Run a Solver

Use [`CreateDefaultSolver`](create.go) to create a solver that can run multiple metaheuristics in parallel, exchanging information (solutions) between them:

```go
import (
    "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/logger/channel"
)

env := MyEnvironment{} // Your implementation
mh := []rko.MetaHeuristic{rko.GA, rko.SA} // Choose metaheuristics
logger := channel.DefaultLogger("MyProblem") // Create a logger

solver := rko.CreateDefaultSolver(mh, env, logger)

result := solver.Solve()
fmt.Printf("Best solution: %+v\n", result)
```

## Documentation

- [`definition.Environment`](definition/definition.go) - The interface you must implement for your optimization problem.
- [`logger.Logger`](logger/definition.go) - Implement this to create your own custom logger.
- [`configuration` package](configuration/README.md) - Configure metaheuristics and solver behavior via code or YAML files.
- [`rko.CreateDefaultSolver`](create.go) - Factory for solvers that run multiple metaheuristics in parallel.
- [`logger/channel`](logger/channel/create.go) - Channel logger (detailed logging with real-time progress tracking).
- [`logger/stdout`](logger/stdout/create.go) - Stdout logger (prints directly to standard output).

## Maintainers

Lucas Mendes - [GitHub](https://github.com/lucasmends)

## License

This project is licensed under the Academic Free License v3.0. See [LICENSE](LICENSE)