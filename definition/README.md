# `definition` Package

Core interfaces and types for defining optimization problems and solvers in RKO-Go.

## Overview

RKO-Go optimizes problems by encoding a solution as a **random key**: a fixed-length vector of `float64` values in `[0, 1)`. Every metaheuristic in this library (GA, BRKGA, SA, ILS, VNS, MultiStart, LNS) works on this generic encoding. The only thing you write is a mapping between a random key and your problem: how many keys it needs, what a key costs, and how to turn it into a real answer. That mapping is the `Environment` interface. Implementing it is the only work required to use the library.

`Cost` returns an `int`, and every metaheuristic in this library **minimizes** it. If your problem is naturally a maximization, negate the objective in `Cost`.

Solvers are driven by `context.Context`. `Solve(ctx)` runs until `ctx` is cancelled or its deadline expires; there is no separate stop/time-limit parameter. A time limit is set by calling `Solve` with a context created via `context.WithTimeout` (or `context.WithDeadline`).

## Environment

Implement this interface for your problem:

```go
type Environment interface {
    NumKeys() int
    Cost(r RandomKey) int
    Decode(r RandomKey) any
    SwapSearch() [][2]int
}
```

| Method | Description |
|---|---|
| `NumKeys() int` | Number of keys (variables) in the problem. |
| `Cost(r RandomKey) int` | Objective value for a random key. Lower is better. |
| `Decode(r RandomKey) any` | Converts a random key into your problem's representation. |
| `SwapSearch() [][2]int` | `(start, end)` index pairs the swap local search iterates over, trying to swap every pair `i, j` with `start <= i < j < end`. |

### Example

```go
package myproblem

import "github.com/RKO-solver/rko-go/definition"

type MyEnv struct{}

func (e MyEnv) NumKeys() int {
    return 10
}

func (e MyEnv) Cost(r definition.RandomKey) int {
    order := r.SortedIndex()
    cost := 0
    for i := 1; i < len(order); i++ {
        cost += order[i] - order[i-1] // replace with your objective function
    }
    return cost
}

func (e MyEnv) Decode(r definition.RandomKey) any {
    return r.SortedIndex()
}

func (e MyEnv) SwapSearch() [][2]int {
    return [][2]int{{0, e.NumKeys()}}
}
```

`SortedIndex` is the standard way to decode a random key for permutation-style problems (see the paper linked below). For other problem types, decode `r` however your problem needs.

## RandomKey

```go
type RandomKey []float64
```

| Method | Description |
|---|---|
| `SortedIndex() []int` | Indices that would sort the key ascending; the standard permutation decoding. |
| `Len() int` | Number of elements in the key. |
| `Clone() RandomKey` | Independent copy of the key. |
| `Equals(other RandomKey) bool` | Element-wise equality. Exact float comparison, no tolerance. |
| `Similarity(key RandomKey) float64` | Cosine similarity with another key of the same length. |

## Solver

Implemented by every metaheuristic (GA, SA, ILS, VNS, MultiStart, BRKGA, LNS). You only need this interface if you are writing a custom solver.

```go
type Solver interface {
    Solve(ctx context.Context) Result
    Name() string
    SetRG(rg *random.Generator)
    Print()
}
```

| Method | Description |
|---|---|
| `Solve(ctx context.Context) Result` | Runs the metaheuristic until `ctx` is done, returns the best `Result` found. |
| `Name() string` | Solver name, used in logs. |
| `SetRG(rg *random.Generator)` | Injects the shared random number generator. |
| `Print()` | Prints the solver's configuration to stdout. |

## Result

```go
type Result struct {
    Solution        any
    Cost            int
    TimeSpentSecond float64
}
```

`Solution` is `Decode(bestKey)`, `Cost` is the objective value of that key, and `TimeSpentSecond` is the wall-clock time the solver spent.

## Context helpers

`WithStartTime` and `StartTime` let a solver measure elapsed time against a shared start, even when the solver did not start the clock itself.

```go
func WithStartTime(ctx context.Context, t time.Time) context.Context
func StartTime(ctx context.Context) time.Time
```

`StartTime` reads the value stored by `WithStartTime`; if none was stored, it returns `time.Now()`. A custom `Solve` implementation typically starts with:

```go
func (s *MySolver) Solve(ctx context.Context) definition.Result {
    start := definition.StartTime(ctx)
    for ctx.Err() == nil {
        // iterate, checking ctx.Err() to stop on cancellation/deadline
    }
    elapsed := time.Since(start).Seconds()
    // ...
}
```

## Files
- `definition.go`: `Environment`, `Solver`, `Result`, `RandomKey`, and the context start-time helpers.

## See Also
- [Project README](../README.md)
- [logger package](../logger/README.md)
- [configuration package](../configuration/README.md)

---
This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
