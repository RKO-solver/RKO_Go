# `problem` Package

The interfaces and types that connect your problem to RKO-Go: `Environment` (your problem), `Cloner` (a copy of it per metaheuristic) and `RandomKey` (a solution).

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

The metaheuristics run in parallel and call the `Environment` methods concurrently. If your `Environment` needs a strategy to handle those concurrent calls, implementing it is up to you. To give each metaheuristic a dedicated `Environment`, implement `Cloner`:

```go
type Cloner interface {
    Clone() Environment
}
```

`rko.New` and `config.NewSolver` then give each metaheuristic its own `Clone()`. `rko.NewFromMetaheuristics` does not: when you build the metaheuristics yourself, clone the environment for each one (see [Building a solver by hand](../metaheuristic/README.md#building-a-solver-by-hand)).

For ordering problems, decode with `r.SortedIndex()`: the indices of `r` in ascending order of their keys, a permutation. The [root README](../README.md#quick-start) solves a travelling salesman problem this way.

Other problems decode keys however they need. This 0-1 knapsack selects item `i` when its key is at least `0.5`, and negates the value because the library minimizes:

```go
package knapsack

import "github.com/RKO-solver/rko-go/problem"

// Knapsack selects items to maximize value without exceeding Capacity.
type Knapsack struct {
	Values, Weights []int
	Capacity        int
}

// NumKeys returns one key per item.
func (k Knapsack) NumKeys() int { return len(k.Values) }

// Decode selects item i when its key is at least 0.5.
func (k Knapsack) Decode(r problem.RandomKey) any {
	items := make([]int, 0)
	for i, key := range r {
		if key >= 0.5 {
			items = append(items, i)
		}
	}
	return items
}

// Cost returns the negated value of the selection, or the excess weight when it is over capacity.
func (k Knapsack) Cost(r problem.RandomKey) int {
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

## See also

- [Project README](../README.md)
- [`config` package](../config/README.md)
- [`metaheuristic` package](../metaheuristic/README.md)
- [`logger` package](../logger/README.md)
