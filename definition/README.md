# `definition` Package

The interfaces and types that connect your problem to RKO-Go: `Environment` (your problem), `RandomKey` (a solution), `Solver` and `Result` (a metaheuristic), and the context helpers that apply the iteration limit.

## Random keys

Every metaheuristic in RKO-Go searches over **random keys**: fixed-length vectors of `float64` values in `[0, 1)`. They never see your problem directly. Your `Environment` tells them how many keys a solution has, what a vector of keys costs, and how to turn it into an answer.

## Environment

```go
type Environment interface {
    NumKeys() int
    Cost(r RandomKey) int
    Decode(r RandomKey) any
    SwapSearch() [][2]int
    ProblemName() string
}
```

| Method | Returns |
|---|---|
| `NumKeys()` | The number of keys in a solution. |
| `Cost(r)` | The objective value of `r`. **Every metaheuristic minimizes it**; negate it to maximize. |
| `Decode(r)` | `r` in your problem's representation. `Solve` returns `Decode` of the best solution. |
| `SwapSearch()` | `{start, end}` ranges for the swap local search, which tries swapping keys `i` and `j` for every `start <= i < j < end`. Return `[][2]int{{0, NumKeys()}}` to allow any pair, or `nil` to disable it. |
| `ProblemName()` | The problem name. It labels the logger and its output files, such as the CSV and history files, so keep it short and safe in file names. |

The metaheuristics run in parallel and call `Cost` and `Decode` concurrently, so these methods must not modify shared state.

For ordering problems, decode with `r.SortedIndex()`: the indices of `r` in ascending order of their keys, a permutation. The [root README](../README.md#quick-start) solves a travelling salesman problem this way.

Other problems decode keys however they need. This 0-1 knapsack selects item `i` when its key is at least `0.5`, and negates the value because the library minimizes:

```go
package knapsack

import "github.com/RKO-solver/rko-go/definition"

// Knapsack selects items to maximise value without exceeding Capacity.
type Knapsack struct {
	Values, Weights []int
	Capacity        int
}

// NumKeys returns one key per item.
func (k Knapsack) NumKeys() int { return len(k.Values) }

// Decode selects item i when its key is at least 0.5.
func (k Knapsack) Decode(r definition.RandomKey) any {
	items := make([]int, 0)
	for i, key := range r {
		if key >= 0.5 {
			items = append(items, i)
		}
	}
	return items
}

// Cost returns the negated value of the selection, or the excess weight when it is over capacity.
func (k Knapsack) Cost(r definition.RandomKey) int {
	value, weight := 0, 0
	for _, i := range k.Decode(r).([]int) {
		value += k.Values[i]
		weight += k.Weights[i]
	}
	if weight > k.Capacity {
		return weight - k.Capacity
	}
	return -value
}

// SwapSearch disables the swap local search.
func (k Knapsack) SwapSearch() [][2]int { return nil }

// ProblemName returns the name that labels the logger and its output files.
func (k Knapsack) ProblemName() string { return "knapsack" }
```

## RandomKey

```go
type RandomKey []float64
```

| Method | Returns |
|---|---|
| `SortedIndex() []int` | The indices that sort the keys ascending (the standard permutation decoding). |
| `Len() int` | The number of keys. |
| `Clone() RandomKey` | An independent copy. |
| `Equals(other RandomKey) bool` | Whether both have the same length and exactly the same values. |
| `Similarity(other RandomKey) float64` | The cosine similarity with a key of the same length. |

## Solver and Result

Every metaheuristic implements `Solver`. You only implement it to add your own metaheuristic.

```go
type Solver interface {
    Solve(ctx context.Context) Result
    Name() string
    Print()
}

type Result struct {
    Solution any           // Decode of the solver's solution
    Cost     int           // its cost
    Elapsed  time.Duration // from the start of the run until Solve returned
}
```

| Method | Does |
|---|---|
| `Solve(ctx)` | Searches until `ctx` is done and returns a `Result`. |
| `Name()` | Returns the name used in logs. |
| `Print()` | Prints the solver's configuration to stdout. |

The run's best solution is the one in the shared solution pool, which is what `rko.Solver.Solve` returns.

## Context helpers

`rko.Solver.Solve` stores the run's start time and iteration limit in the context it passes to every metaheuristic. These helpers read them:

| Function | Does |
|---|---|
| `WithStartTime(ctx, t) context.Context` | Returns a context carrying the run's start time. |
| `StartTime(ctx) time.Time` | Returns the stored start time, or `time.Now()` if there is none. |
| `WithMaxIterations(ctx, n) context.Context` | Returns a context carrying the iteration limit of each metaheuristic. Zero or negative means no limit. |
| `GetMHContext(ctx) (context.Context, *IterationCounter)` | Returns a context for one metaheuristic and the counter that cancels it at the limit. Without a limit, it returns `ctx` itself. |

`IterationCounter` methods:

| Method | Does |
|---|---|
| `Step()` | Counts one iteration and cancels the metaheuristic's context when the limit is reached. |
| `Iteration() int` | Returns the number of completed iterations. |
| `Stop()` | Releases the context. Call it with `defer`. |

## Writing a metaheuristic

Call `GetMHContext` first, count one iteration of your outer loop with `Step`, and stop when `ctx.Err()` is not nil. That single check covers the time limit, the iteration limit and cancellation. Offer improvements to the shared pool so the other metaheuristics can use them; it keeps only solutions that beat its best.

```go
package mysolver

import (
	"context"
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// RandomSearch samples random keys and keeps the best one.
type RandomSearch struct {
	env  definition.Environment
	pool *solution.Pool
	rg   *random.Generator
	log  logger.SolverLogger // from Logger.SolverLogger("RandomSearch")
}

// NewRandomSearch receives its dependencies, including the generator, through the constructor.
func NewRandomSearch(env definition.Environment, pool *solution.Pool, rg *random.Generator, log logger.SolverLogger) *RandomSearch {
	return &RandomSearch{env: env, pool: pool, rg: rg, log: log}
}

func (s *RandomSearch) Name() string { return "RandomSearch" }
func (s *RandomSearch) Print()       { fmt.Println(s.Name()) }

func (s *RandomSearch) Solve(ctx context.Context) definition.Result {
	ctx, counter := definition.GetMHContext(ctx)
	defer counter.Stop()
	start := definition.StartTime(ctx)

	key := rk.Generate(s.env, s.rg)
	best := &metaheuristc.RandomKeyValue{RK: key, Cost: s.env.Cost(key)}

	for ; ctx.Err() == nil; counter.Step() {
		key := rk.Generate(s.env, s.rg)
		if cost := s.env.Cost(key); cost < best.Cost {
			best = &metaheuristc.RandomKeyValue{RK: key, Cost: cost}
			s.pool.Add(best, start)
		}
		s.log.Register(best.Cost, best.Cost, counter.Iteration(), time.Since(start))
	}

	return definition.Result{
		Solution: s.env.Decode(best.RK),
		Cost:     best.Cost,
		Elapsed:  time.Since(start),
	}
}
```

Pass it to `rko.CreateFullSolver` with the other solvers; the [`metaheuristc` package](../metaheuristc/README.md#building-a-solver-by-hand) shows how. Report progress once per iteration with `Register` and occasional events with `Event`; see the [`logger` package](../logger/README.md).

## See also

- [Project README](../README.md)
- [`config` package](../config/README.md)
- [`metaheuristc` package](../metaheuristc/README.md)
- [`logger` package](../logger/README.md)
