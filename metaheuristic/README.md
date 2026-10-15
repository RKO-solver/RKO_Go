# `metaheuristic` Package

The metaheuristic implementations and the pieces they share: local searches, the solution pool and random-key helpers.

Most programs never import these packages: [`rko.New`](../README.md#quick-start) and the [`config` package](../config/README.md) build everything for you. Come here to build a solver by hand, combine local searches yourself, or add your own metaheuristic next to the built-in ones.

| Package | Contents |
|---|---|
| `metaheuristic/multistart`, `sa`, `ga`, `brkga`, `ils`, `vns`, `lns` | The seven metaheuristics. |
| `metaheuristic/search` | Local searches. |
| `metaheuristic/solution` | The solution pool shared by every metaheuristic. |
| `metaheuristic/rk` | `Generate`, `Reset`, `Copy` and `Shake` for random keys. |
| [`metaheuristic/tuning`](tuning/README.md) | The Q-learning tuner that chooses parameter values while a metaheuristic runs. |
| `metaheuristic` | `Type` (`ParseType`, `String`, `Short`, `Implemented`) to name a metaheuristic; `RandomKeyValue` (a key and its cost) and `NewRandomKeyValue`, `List[T]` (a parameter's values), `Sort`, the `Worker` used by `rko.Solver`; `Algorithm` and `Result` for [your own metaheuristic](#writing-a-metaheuristic); `LocalSearch` and `Print`; and the [context helpers](#context-helpers). |

## Metaheuristics

| Metaheuristic | Type | `metaheuristic.Type` constant | One iteration is |
|---|---|---|---|
| MultiStart | `*multistart.MultiStart` | `metaheuristic.MultiStart` | a random solution improved by the local search |
| Simulated Annealing | `*sa.SimulatedAnnealing` | `metaheuristic.SA` | one temperature step of `Iterations` moves |
| Genetic Algorithm | `*ga.GA` | `metaheuristic.GA` | one generation |
| BRKGA | `*brkga.BRKGA` | `metaheuristic.BRKGA` | one generation |
| Iterated Local Search | `*ils.ILS` | `metaheuristic.ILS` | one shake followed by the local search |
| Variable Neighbourhood Search | `*vns.VNS` | `metaheuristic.VNS` | one pass through neighbourhoods 1 to 3, restarting at 1 after each improvement |
| Large Neighbourhood Search | `*lns.LNS` | `metaheuristic.LNS` | one rebuild of some keys followed by the local search |

The iteration limit counts these iterations (see [Stopping](#stopping)). Every metaheuristic offers its improvements to the shared pool.

`metaheuristic.Type` also declares `GRASP`, `VLNS`, `ALNS` and `IPR`, which are not implemented: `rko.New` and `config.NewSolver` reject them. `t.Implemented()` reports whether `rko.New` can build `t`. `t.String()` returns the full name (`"Simulated Annealing"`) and `t.Short()` the abbreviation (`"SA"`; `"MS"` for MultiStart). `ParseType(label)` accepts `"MultiStart"`, `"SA"`, `"GA"`, `"VNS"`, `"ILS"`, `"BRKGA"`, `"GRASP"`, `"LNS"`, `"VLNS"`, `"ALNS"` and `"IPR"` (case-insensitive) and returns `-1` for anything else.

- **MultiStart** draws a random solution and improves it with the local search, over and over.
- **SA** starts from the pool's best solution. Each move shakes the current solution and runs the local search on it; a worse result is accepted with probability `exp(-Δ / T)`. It reheats up to `QtdReheat` times. When the temperature reaches `TemperatureGoal`, SA ends, or, with `Restart`, starts over: the temperature goes back to `TemperatureInitial`, the current solution becomes a new random one and the `QtdReheat` reheats are available again. With `Restart`, SA runs until its context stops.
- **GA** picks parents by 3-way tournament, applies uniform crossover with mutation, and runs the local search on one random child per generation.
- **BRKGA** keeps the elite, breeds each child from an elite and a non-elite parent, adds random mutants, and runs the local search on one random individual per generation.
- GA and BRKGA restart their population after more than `MaxGenerationNoImprovement` generations without improvement, keeping the best solution.
- **ILS** shakes the current solution and runs the local search, keeping the result when it beats the best. After more than 10 (and fewer than 1000) consecutive iterations without improvement, the shake range grows by 5% per iteration, falling back to 1.1 times the configured range if it passes `1.0`; an improvement restores it. With `MetropolisCriterion` and a deadline, it may also continue from a worse solution.
- **VNS** shakes with a growing intensity: neighbourhood `k` shakes with `r` drawn from `[k × Rate, (k+1) × Rate)` (a shake applies `int(NumKeys × r) + 1` random moves), then runs the local search.
- **LNS** picks some keys of the best solution at random and, for each one, tries a random value in each interval of a Farey sequence and keeps the best; then it runs the local search.

### Logged events

Every metaheuristic calls `Register` once per iteration (see the [`logger` package](../logger/README.md#writing-a-metaheuristic)). These also call `Event`, which is shown only at `LevelVerbose`:

| Metaheuristic | Event key | Value | Fires |
|---|---|---|---|
| GA, BRKGA | `stalled` | group: `generation` (int), `no-improvement` (int) | After each generation that did not improve the best solution. `no-improvement` counts consecutive such generations. |
| GA, BRKGA | `restart` | string `limit generations without improvement` | When the population restarts. |
| SA | `temperature` | float | After every temperature step, with the new temperature. |
| SA | `reheat` | int | At each reheat, with the reheat number (`1` for the first). |
| SA | `restart` | group: `iteration` (int) | At each restart (with `Restart`), with the outer-loop iteration. From then on, the best cost SA passes to `Register` is the best of the new run, so it can go up. |
| VNS | `improved` | group: `iteration` (int), `local` (int) | When a neighbour beats the best solution. |

MultiStart, ILS and LNS emit no events. A group prints as `stalled=[generation=152 no-improvement=3]`.

Parameters, their meaning and their defaults are in the [`config` package](../config/README.md#parameters-section).

## Building a solver by hand

Create one pool, then each metaheuristic with its configuration and local search, and pass them to `rko.NewFromMetaheuristics`. It returns `rko.ErrNoMetaheuristics` if the slice is empty (`rko.WithMetaheuristics` is ignored there):

```go
rg := random.New(42)
log := printer.Default("knapsack")
defer log.Close()
pool := solution.NewPool(env, rg, log)

ilsConfig := ils.DefaultConfig()
ilsConfig.MetropolisCriterion = new(true)

brkgaConfig := brkga.DefaultConfig()
brkgaConfig.PopulationSize = []int{100}

brkgaRG := rg.Split()
metaheuristics := []metaheuristic.Algorithm{
	ils.New(env, rg.Split(), pool, log,
		ils.WithConfig(ilsConfig),
		ils.WithSearch(search.NewMirror(env))),
	brkga.New(env, brkgaRG, pool, log,
		brkga.WithConfig(brkgaConfig),
		brkga.WithSearch(search.NewRVND(env, pool, brkgaRG, []search.Type{search.Mirror, search.Farey}))),
}

solver, err := rko.NewFromMetaheuristics(env, log, pool, metaheuristics,
	rko.WithTimeLimit(30*time.Second), rko.WithMaxIterations(1000))
if err != nil {
	fmt.Println(err)
	return
}
items := solver.Solve(context.Background()).([]int)
```

- Give every metaheuristic the pool and logger you pass to `NewFromMetaheuristics`: `Solve` returns the best solution of that pool.
- Give each metaheuristic its own generator with `rg.Split()`: a `random.Generator` is not safe for concurrent use. A local search that draws random numbers (`Farey`, `NelderMead`, RVND) should use the same generator as its metaheuristic.
- If you need a dedicated `Environment` per metaheuristic (to handle concurrency), give each metaheuristic and its local search its own `env.Clone()` (see [`problem.Cloner`](../problem/README.md#environment)).
- `rko.WithTimeLimit(d)` and `rko.WithMaxIterations(n)` set the limits; a zero or negative value means no limit, and with no option there is none.
- Your own [`metaheuristic.Algorithm`](#writing-a-metaheuristic) can go in the same slice.

### Constructors

Each package has one constructor, `New`, that returns the concrete type (`*ils.ILS`, ...):

```go
New(env, rg, pool, log, opts...)
```

`env` is a `problem.Environment`, `rg` a `*random.Generator`, `pool` a `*solution.Pool` and `log` a `logger.Logger`. Without options it uses the package's `DefaultConfig()` and the default local search, `search.NewDefault` (RVND over all four searches). The options:

| Option | Sets |
|---|---|
| `WithConfig(cfg)` | The parameters. Not available in `multistart`. |
| `WithSearch(local)` | A ready `metaheuristic.LocalSearch`. |
| `WithSearchType(t)` | The local search, built with `search.New(t, ...)`. |

- If `WithSearch` and `WithSearchType` are both given, the last one wins.
- A nil config or local search means the default.

### Configuration structs

`sa`, `ga`, `brkga`, `ils`, `vns` and `lns` each have a `Config` struct with exactly the fields of the YAML sections, with the same names and element types; `DefaultConfig()` returns it at its defaults. See the [`config` package](../config/README.md#parameters-section) for each field. Every field is a `metaheuristic.List[T]` (a `[]T`) of the values to use; `DefaultConfig()` gives each one a single value. The exceptions are `sa.Config.Restart` and `ils.Config.MetropolisCriterion`, plain `*bool` fields set with `new(true)` or `new(false)`. Fields are used as is: change them on a default configuration rather than starting from a zero struct, keep every list non-empty and never set those two to nil.

```go
cfg := brkga.DefaultConfig()
cfg.PopulationSize = []int{100}                // one value: fixed
cfg.EliteRatio = []float64{0.1, 0.2}           // several values: tuned
```

## Parameter tuning

When any tunable field of a `Config` has more than one value, `New` creates a Q-learning tuner (`tuning.NewQLearning`, from the [`tuning` package](tuning/README.md)) over every combination of the values. Otherwise there is no tuner: the first value of each list is used, and a fixed seed gives the same results. Q-learning is the only tuner for now.

Each cycle, the tuner chooses the configuration; after the cycle it learns from a reward (higher is better):

| Metaheuristic | Tunable fields | Cycle | Reward |
|---|---|---|---|
| BRKGA | `PopulationSize`, `EliteRatio`, `MutantRatio`, `CrossoverAlpha`, `MutationAlpha`, `MaxGenerationNoImprovement` | generation | `1 + 1/PopulationSize` if the best improved, else `(best - best of the generation) / abs(best)` |
| GA | `PopulationSize`, `CrossoverAlpha`, `MutationAlpha`, `MaxGenerationNoImprovement` | generation | same as BRKGA |
| SA | `Alpha`, `TemperatureInitial`, `TemperatureGoal`, `TemperatureReheat`, `ShakeMin`, `ShakeMax`, `QtdReheat`, `Iterations` | temperature step | `1` if the best improved, else `(best - current) / abs(best)` |
| ILS | `ShakeMin`, `ShakeMax` | iteration | same as SA |
| VNS | `Rate` | outer iteration | `1` if the best improved, else `(best - last neighbour) / abs(best)` |
| LNS | `BetaMin`, `BetaMax` | iteration | same as SA |

`abs(best)` is taken as `1` when the best cost is `0`.

SA's `Restart` and ILS's `MetropolisCriterion` are plain bools, not lists, and are never tuned.

- GA and BRKGA allocate the population at the largest `PopulationSize`; when the size grows, new random individuals are added.
- SA starts at the first `TemperatureInitial`; a tuned `TemperatureInitial` applies at the reheats and the restarts.
- Only `config.NewSolver` checks that every `TemperatureGoal` is below every `TemperatureInitial`. `sa.New` does not: with a hand-built `Config` whose goal is at or above the initial temperature, SA with `Restart` restarts at every temperature step.
- The tuner explores less as the run advances; see [Progress](#progress).
- `New` panics if BRKGA has an `EliteRatio` and `PopulationSize` combination whose elite size, `int(EliteRatio * PopulationSize)`, is outside `1` to `PopulationSize - 1`, or if LNS has a `BetaMin` and `BetaMax` combination with `int(BetaMin * NumKeys) >= int(BetaMax * NumKeys)`.

### Progress

Before and after each cycle, the metaheuristic gives the tuner the run's **progress**, a number in `[0, 1]`. It is measured by time when the metaheuristic's context has a deadline, by cycles against `MaxIterations` when only that is set, and by cycles on a curve with no fixed end otherwise:

| Limits | Progress by |
|---|---|
| A time limit (including the `rko.DefaultTimeLimit` of 150 s that `rko.New` sets) | time, up to the time limit |
| A deadline on the `ctx` given to `Solve`, with or without a time limit | time, up to the earlier of the two |
| A time limit and `MaxIterations` | time; `MaxIterations` is not used in the progress, but it may still end the run first |
| `MaxIterations` only, and no deadline on `ctx` | cycles over `MaxIterations` |
| No limit and no deadline on `ctx` (the run continues until `ctx` is cancelled) | cycles, on a curve with no fixed end |

**By time**, progress grows linearly from `0` at the start to `1` at the deadline:

```text
progress = min(1, (now − start) / (deadline − start))
```

`start` is the run's start time (`metaheuristic.StartTime(ctx)`) and `deadline` the context's deadline. If `ctx` has no deadline, or the deadline is not after `start` (for example, it has already passed), the run is treated as having no deadline and progress follows the `MaxIterations` or no-limit rows.

**By cycles over `MaxIterations`**, progress grows linearly from `0` to `1` at the last iteration:

```text
progress = min(1, i / MaxIterations)
```

`i` is the number of cycles completed so far. One cycle is one iteration of the metaheuristic (see [Metaheuristics](#metaheuristics)), the same count that `MaxIterations` limits.

**By cycles with no limit**, the length of the run is unknown, so progress follows a curve with no fixed end:

```text
L = ln(1 + i / States())
progress = L / (1 + L)
```

`i` is the number of cycles completed so far. `States()` is the number of value combinations: the product of the list lengths of the tunable `Config` fields.

This curve applies only to runs with no limit. Progress starts at `0`, rises fast at first and then more and more slowly, and never reaches `1`. Because `i` is divided by `States()`, a larger space keeps exploring for more cycles:

| `i / States()` | Progress |
|---|---|
| 0 | 0.00 |
| 1 | 0.41 |
| e − 1 ≈ 1.72 | 0.50 |
| 10 | 0.71 |
| 100 | 0.82 |
| 1000 | 0.87 |

For example, with 8 combinations, progress is 0.50 after about 14 cycles and 0.82 after 800.

`QLearning` lowers its exploration rate and its learning rate as progress rises; see [Exploration rate](tuning/README.md#exploration-rate) and [Learning rate](tuning/README.md#learning-rate). Its exploration rate settles at its floor of `0.1` at progress `0.9`: at `0.9 × MaxIterations` cycles with `MaxIterations`, but only after about `8100 × States()` cycles with no limit.

### Learned policy

When tuned, the metaheuristic records the learned policy with its metaheuristic logger at the end of the run. After the run, print it with the logger's `PrintPolicies` or save it with `SavePoliciesCSV` (see [Policies](../logger/README.md#policies)). The policy columns `p0`, `p1`, ... are the tunable fields, in the field order of the `Config` struct; a tunable field with a single value still has its column. SA's `Restart` and ILS's `MetropolisCriterion` are not tuned, so they have none:

| Metaheuristic | `p0`, `p1`, ... |
|---|---|
| BRKGA | `PopulationSize`, `EliteRatio`, `MutantRatio`, `CrossoverAlpha`, `MutationAlpha`, `MaxGenerationNoImprovement` |
| GA | `PopulationSize`, `CrossoverAlpha`, `MutationAlpha`, `MaxGenerationNoImprovement` |
| SA | `Alpha`, `TemperatureInitial`, `TemperatureGoal`, `TemperatureReheat`, `ShakeMin`, `ShakeMax`, `QtdReheat`, `Iterations` |
| ILS | `ShakeMin`, `ShakeMax` |
| VNS | `Rate` |
| LNS | `BetaMin`, `BetaMax` |

## Local searches

A local search improves a solution in place, never making it worse, and returns early when `ctx` is done.

| `search.Type` | Tries |
|---|---|
| `Swap` | Swapping every pair of keys inside each `Environment.SwapSearch()` range. |
| `Mirror` | Replacing each key `v` with `1 - v`. |
| `Farey` | For each key, a random value in each interval of a Farey sequence. |
| `NelderMead` | A Nelder-Mead-style search that blends the solution with two solutions from the pool. |
| `RVND` | The searches in its list in random order. A search that fails to improve is dropped; an improvement restores the full list. It ends when none is left. |

| Constructor | Builds |
|---|---|
| `NewSwap(env)` | `Swap` |
| `NewMirror(env)` | `Mirror` |
| `NewFarey(env, rg)` | `Farey` |
| `NewNelderMead(env, pool, rg)` | `NelderMead` |
| `NewRVND(env, pool, rg, types)` | RVND over `types` (`RVND` entries are removed) |
| `NewDefault(env, pool, rg)` | RVND over `Swap`, `Mirror`, `Farey`, `NelderMead` |
| `New(searchType, env, pool, rg)` | The single search for `Swap`, `Mirror`, `Farey` or `NelderMead`; RVND over `Swap`, `Mirror`, `Farey` for `RVND`; `NewDefault` for any other value |

All return a `metaheuristic.LocalSearch`:

```go
type LocalSearch interface {
    Search(ctx context.Context, rko *metaheuristic.RandomKeyValue)
    String() string
}
```

`ParseType(name)` parses `"Swap"`, `"Mirror"`, `"Farey"`, `"Nelder"` or `"RVND"` (case-insensitive; `-1` if unknown). `t.String()` returns the display name (`"Nelder-Mead"` for `NelderMead`). `metaheuristic.Print(name, local)` prints a metaheuristic's name and its local search; use it in `Print` when your metaheuristic has one.

## Solution pool

`solution.Pool` is the thread-safe store through which the metaheuristics share solutions. `rko.Solver.Solve` returns the decoded best solution of the pool.

| Function | Creates |
|---|---|
| `NewPool(env, rg, log)` | A pool seeded with 10 random solutions, keeping at most 10. |
| `NewPoolSize(env, rg, log, maxSize, initialSize)` | A pool seeded with `initialSize` random solutions. A `maxSize` of 0 or less means no cap; a cap below `initialSize` is raised to it. |

| Method | Does |
|---|---|
| `Add(s, start)` | Adds a copy of `s` only if its cost is strictly lower than the pool's best, and reports it to `log.AddSolutionPool` with the time elapsed since `start`. |
| `Best()` | Returns a copy of the best solution, or nil if the pool is empty. |
| `BestCost()` | Returns the best cost, or `math.MaxInt` if the pool is empty. |
| `Solution(i)` | Returns a copy of the `i`-th best solution (`0` is the best). |
| `Len()` | Returns the number of solutions. |

## Stopping

Each metaheuristic calls `metaheuristic.NewSolverContext(ctx)` at the start of `Solve`, counts one iteration of its outer loop with `Step`, and stops when its context is done: at the time limit, at the iteration limit or on cancellation. The limits are set by `rko.WithTimeLimit` and `rko.WithMaxIterations` or the [`config` package](../config/README.md#stopping). See [Context helpers](#context-helpers).

## Result

Every metaheuristic implements the `metaheuristic.Algorithm` interface, and `Solve` returns a `metaheuristic.Result`. You only implement `Algorithm` to add your own metaheuristic.

```go
type Algorithm interface {
    Solve(ctx context.Context) metaheuristic.Result
    Name() string
    Print()
}

type Result struct {
    Solution any           // Decode of the metaheuristic's solution
    Cost     int           // its cost
    Elapsed  time.Duration // from the start of the run until Solve returned
}
```

| Method | Does |
|---|---|
| `Solve(ctx)` | Searches until `ctx` is done and returns a `Result`. |
| `Name()` | Returns the name used in logs. |
| `Print()` | Prints the metaheuristic's configuration to stdout. |

The run's best solution is the one in the shared solution pool, which is what `rko.Solver.Solve` returns.

## Context helpers

`rko.Solver.Solve` stores the run's start time and iteration limit in the context it passes to every metaheuristic. These helpers read them:

| Function | Does |
|---|---|
| `WithStartTime(ctx, timeLimit) (context.Context, context.CancelFunc)` | Returns a context carrying the current time as the run's start time. If `timeLimit > 0`, it also sets a deadline `timeLimit` after that start (an earlier deadline on `ctx` is kept). Call `defer cancel()`. |
| `StartTime(ctx) time.Time` | Returns the stored start time, or `time.Now()` if there is none. |
| `WithMaxIterations(ctx, n) context.Context` | Returns a context carrying the iteration limit of each metaheuristic. Zero or negative means no limit. It is not the option `rko.WithMaxIterations`, which sets the limit on a `Solver`. |
| `NewSolverContext(ctx) (context.Context, *IterationCounter)` | Returns a context for one metaheuristic and the counter that cancels it at the limit. Without a limit, it returns `ctx` itself. |

`IterationCounter` methods:

| Method | Does |
|---|---|
| `Step()` | Counts one iteration and cancels the metaheuristic's context when the limit is reached. |
| `Iteration() int` | Returns the number of completed iterations. |
| `Limit() int` | Returns the iteration limit, or `-1` if there is none. |
| `Stop()` | Releases the context. Call it with `defer`. |

## Writing a metaheuristic

Call `NewSolverContext` first, count one iteration of your outer loop with `Step`, and stop when `ctx.Err()` is not nil. That single check covers the time limit, the iteration limit and cancellation. Offer improvements to the shared pool so the other metaheuristics can use them; it keeps only solutions that beat its best.

```go
package randomsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

// RandomSearch samples random keys and keeps the best one.
type RandomSearch struct {
	env  problem.Environment
	pool *solution.Pool
	rg   *random.Generator
	log  logger.MetaheuristicLogger // from Logger.MetaheuristicLogger("RandomSearch")
}

var _ metaheuristic.Algorithm = (*RandomSearch)(nil)

// NewRandomSearch receives its dependencies, including the generator, through the constructor.
func NewRandomSearch(env problem.Environment, pool *solution.Pool, rg *random.Generator, log logger.MetaheuristicLogger) *RandomSearch {
	return &RandomSearch{env: env, pool: pool, rg: rg, log: log}
}

func (s *RandomSearch) Name() string { return "RandomSearch" }
func (s *RandomSearch) Print()       { fmt.Println(s.Name()) }

func (s *RandomSearch) Solve(ctx context.Context) metaheuristic.Result {
	ctx, counter := metaheuristic.NewSolverContext(ctx)
	defer counter.Stop()
	start := metaheuristic.StartTime(ctx)

	key := rk.Generate(s.env, s.rg)
	best := &metaheuristic.RandomKeyValue{RK: key, Cost: s.env.Cost(key)}

	for ; ctx.Err() == nil; counter.Step() {
		key := rk.Generate(s.env, s.rg)
		if cost := s.env.Cost(key); cost < best.Cost {
			best = &metaheuristic.RandomKeyValue{RK: key, Cost: cost}
			s.pool.Add(best, start)
		}
		s.log.Register(best.Cost, best.Cost, counter.Iteration(), time.Since(start))
	}

	return metaheuristic.Result{
		Solution: s.env.Decode(best.RK),
		Cost:     best.Cost,
		Elapsed:  time.Since(start),
	}
}
```

Pass it to `rko.NewFromMetaheuristics` with the other metaheuristics, as shown in [Building a solver by hand](#building-a-solver-by-hand). If your metaheuristic has a local search, `Print` can call `metaheuristic.Print(name, local)`. Report progress once per iteration with `Register` and occasional events with `Event`; see the [`logger` package](../logger/README.md).

## See also

- [Project README](../README.md)
- [`config` package](../config/README.md)
- [`problem` package](../problem/README.md)
- [`logger` package](../logger/README.md)
