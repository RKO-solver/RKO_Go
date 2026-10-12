# `metaheuristc` Package

The metaheuristic implementations and the pieces they share: local searches, the solution pool and random-key helpers.

Most programs never import these packages: [`rko.New`](../README.md#quick-start) and the [`config` package](../config/README.md) build everything for you. Come here to build solvers by hand, combine local searches yourself, or add your own metaheuristic next to the built-in ones.

| Package | Contents |
|---|---|
| `metaheuristc/multistart`, `sa`, `ga`, `brkga`, `ils`, `vns`, `lns` | The seven metaheuristics. |
| `metaheuristc/search` | Local searches. |
| `metaheuristc/solution` | The solution pool shared by every metaheuristic. |
| `metaheuristc/rk` | `Generate`, `Reset` and `Shake` for random keys. |
| `metaheuristc/constants` | Default parameter values and `DefaultTimeLimit` (150 s). |
| `metaheuristc` | `RandomKeyValue` (a key and its cost), `Sort` and the `Worker` used by `rko.Solver`. |

## Metaheuristics

| Metaheuristic | Type | `rko` constant | One iteration is |
|---|---|---|---|
| MultiStart | `*multistart.MultiStart` | `MultiStart` | a random solution improved by the local search |
| Simulated Annealing | `*sa.SimulatedAnnealing` | `SA` | one temperature step of `Iterations` moves |
| Genetic Algorithm | `*ga.GA` | `GA` | one generation |
| BRKGA | `*brkga.BRKGA` | `BRKGA` | one generation |
| Iterated Local Search | `*ils.ILS` | `ILS` | one shake followed by the local search |
| Variable Neighbourhood Search | `*vns.VNS` | `VNS` | one pass through neighbourhoods 1 to 3, restarting at 1 after each improvement |
| Large Neighbourhood Search | `*lns.LNS` | `LNS` | one rebuild of some keys followed by the local search |

The iteration limit counts these iterations (see [Stopping](#stopping)). Every metaheuristic offers its improvements to the shared pool.

- **MultiStart** draws a random solution and improves it with the local search, over and over.
- **SA** starts from the pool's best solution. Each move shakes the current solution and runs the local search on it; a worse result is accepted with probability `exp(-Δ / T)`. It reheats up to `QtdReheat` times and ends when the temperature reaches `TemperatureGoal`.
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
| SA | `reheat` | int | At each reheat, with the number of earlier reheats (`0` for the first). |
| VNS | `improved` | group: `iteration` (int), `local` (int) | When a neighbour beats the best solution. |

MultiStart, ILS and LNS emit no events. A group prints as `stalled=[generation=152 no-improvement=3]`.

Parameters, their meaning and their defaults are in the [`config` package](../config/README.md#parameters-section).

## Building a solver by hand

Create one pool, then each metaheuristic with its configuration and local search, and pass them to `rko.NewFromSolvers`:

```go
rg := random.New(42)
log := printer.Default("knapsack")
defer log.Close()
pool := solution.NewPool(env, rg, log)

ilsConfig := ils.DefaultConfig()
ilsConfig.MetropolisCriterion = true

brkgaConfig := brkga.DefaultConfig()
brkgaConfig.PopulationSize = 100

brkgaRG := rg.Split()
solvers := []definition.Solver{
	ils.New(env, rg.Split(), pool, log,
		ils.WithConfig(ilsConfig),
		ils.WithSearch(search.NewMirror(env))),
	brkga.New(env, brkgaRG, pool, log,
		brkga.WithConfig(brkgaConfig),
		brkga.WithSearch(search.NewRVND(env, pool, brkgaRG, []search.Type{search.Mirror, search.Farey}))),
}

solver := rko.NewFromSolvers(env, log, pool, solvers,
	rko.WithTimeLimit(30*time.Second), rko.WithMaxIterations(1000))
items := solver.Solve(context.Background()).([]int)
```

- Give every solver the pool and logger you pass to `NewFromSolvers`: `Solve` returns the best solution of that pool.
- Give each solver its own generator with `rg.Split()`: a `random.Generator` is not safe for concurrent use. A local search that draws random numbers (`Farey`, `NelderMead`, RVND) should use the same generator as its metaheuristic.
- `rko.WithTimeLimit(d)` and `rko.WithMaxIterations(n)` set the limits; a zero or negative value means no limit, and with no option there is none.
- Your own [`definition.Solver`](../definition/README.md#writing-a-metaheuristic) can go in the same slice.

### Constructors

Each package has one constructor, `New`, that returns the concrete type (`*ils.ILS`, ...):

```go
New(env, rg, pool, log, opts...)
```

`env` is a `definition.Environment`, `rg` a `*random.Generator`, `pool` a `*solution.Pool` and `log` a `logger.Logger`. Without options it uses the package's `DefaultConfig()` and the default local search, `search.NewDefault` (RVND over all four searches). The options:

| Option | Sets |
|---|---|
| `WithConfig(cfg)` | The parameters. Not available in `multistart`. |
| `WithSearch(local)` | A ready `search.Local`. |
| `WithSearchType(t)` | The local search, built with `search.New(t, ...)`. |

- If `WithSearch` and `WithSearchType` are both given, the last one wins.
- A nil config or local search means the default.

### Configuration structs

`sa`, `ga`, `brkga`, `ils`, `vns` and `lns` each have a `Config` struct with exactly the fields of the YAML sections, with the same names and types; `DefaultConfig()` returns it at its defaults. See the [`config` package](../config/README.md#parameters-section) for each field. Fields are used as is: change them on a default configuration rather than starting from a zero struct.

## Local searches

A local search improves a solution in place, never making it worse, and returns early when `ctx` is done.

| `search.Type` | Tries |
|---|---|
| `Swap` | Swapping every pair of keys inside each `Environment.SwapSearch()` range. |
| `Mirror` | Replacing each key `v` with `1 - v`. |
| `Farey` | For each key, a random value in each interval of a Farey sequence. |
| `NelderMead` | A Nelder-Mead-style search that blends the solution with two solutions from the pool. |
| `RVND` | The searches in its list in random order, dropping each one that fails to improve, until none is left. |

| Constructor | Builds |
|---|---|
| `NewSwap(env)` | `Swap` |
| `NewMirror(env)` | `Mirror` |
| `NewFarey(env, rg)` | `Farey` |
| `NewNelderMead(env, pool, rg)` | `NelderMead` |
| `NewRVND(env, pool, rg, types)` | RVND over `types` (`RVND` entries are removed) |
| `NewDefault(env, pool, rg)` | RVND over `Swap`, `Mirror`, `Farey`, `NelderMead` |
| `New(searchType, env, pool, rg)` | The single search for `Swap`, `Mirror`, `Farey` or `NelderMead`; RVND over `Swap`, `Mirror`, `Farey` for `RVND`; `NewDefault` for any other value |

All return a `search.Local`:

```go
type Local interface {
    Search(ctx context.Context, rko *metaheuristc.RandomKeyValue)
    String() string
}
```

`ParseType(name)` parses `"Swap"`, `"Mirror"`, `"Farey"`, `"Nelder"` or `"RVND"` (case-insensitive; `-1` if unknown). `t.String()` returns the display name (`"Nelder-Mead"` for `NelderMead`).

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

Each metaheuristic calls `definition.NewSolverContext(ctx)` at the start of `Solve`, counts one iteration of its outer loop with `Step`, and stops when its context is done: at the time limit, at the iteration limit or on cancellation. The limits are set by `rko.WithTimeLimit` and `rko.WithMaxIterations` or the [`config` package](../config/README.md#stopping). See the [`definition` package](../definition/README.md#context-helpers) for the helpers.

## See also

- [Project README](../README.md)
- [`config` package](../config/README.md)
- [`definition` package](../definition/README.md)
- [`logger` package](../logger/README.md)
