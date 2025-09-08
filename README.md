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

You can track the progress of your optimization by providing a custom logger that implements the [`logger.Interface`](logger/definition.go):

```go
type Interface interface {
    Report(bestSolutionCost, localSolutionCost int, elapsed float64)
    Verbose(message string)
    Debug(message string)
    Info(message string)
    SetIdWorker(idWorker int)
    CreateLogger(method string) Interface
    Save()
    SaveFileName(fileName string)
}
```

### Built-in Loggers

- [`logger/basic`](logger/basic/implementation.go): Prints progress and messages to the screen.
- [`logger/csv`](logger/csv/implementation.go): Saves progress and results to a CSV file.

You can implement your own logger by satisfying the [`logger.Interface`](logger/definition.go).

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
}
```

- `NumKeys()` returns the number of keys (variables) in your problem.
- `Cost(r RandomKey)` evaluates the cost (objective function) for a solution.
- `Decode(r RandomKey)` converts a solution from random keys to your problem's representation.

### 3. Create and Run a Solver

Use [`CreateDefaultSolver`](create.go) to create a solver that can run multiple metaheuristics in parallel, exchanging information (solutions) between them:

```go
import (
    "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/logger/basic"
)

env := MyEnvironment{} // Your implementation
mh := []rko.MetaHeuristic{rko.GA, rko.SA} // Choose metaheuristics

solver := rko.CreateDefaultSolver(
    mh,
    env,
    rko.INFO,           // Log level
    false,              // Save report
    basic.CreateLogger(), // Logger implementation
)

result := solver.Solve()
fmt.Printf("Best solution: %+v\n", result)
```

## Documentation

- [`definition.Environment`](definition/definition.go) - The interface you must implement.
- [`logger.Interface`](logger/definition.go) - Implement this to create your own logger.
- [`rko.CreateDefaultSolver`](create.go) - Factory for solvers that run multiple metaheuristics in parallel.
- [`logger/basic`](logger/basic/implementation.go) - Basic logger (prints to screen).
- [`logger/csv`](logger/csv/implementation.go) - CSV logger (saves to file).

## Maintainers

Lucas Mendes - [GitHub](https://github.com/lucasmends)

## License

This project is licensed under the Academic Free License v3.0. See [LICENSE](LICENSE)