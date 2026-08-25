# `metaheuristc` Package

This package holds the implementations of every metaheuristic in RKO-Go: `multistart`, `sa`, `ga` (GA and BRKGA), `ils`, `vns`, `lns`, plus the shared `search` (local search), `solution` (solution pool), `rk` (random-key helpers) and `constants` (default values) subpackages.

Most users never import `metaheuristc` directly — the root [`rko`](../create.go) package or [`configuration`](../configuration/README.md) build and run these solvers for you. Come here when you need to hand-build a solver, pick or combine local searches yourself, or read the exact default value of a tuning parameter.

Everything below is read from source: `metaheuristc/constants/defaults.go`, `metaheuristc/*/definition.go`, `metaheuristc/*/create.go`, `metaheuristc/*/interface.go`, and `metaheuristc/search/*.go`.

## The seven metaheuristics

| Name | Subpackage | Solver type | Configuration type | Default-config constructor |
|---|---|---|---|---|
| MultiStart | `metaheuristc/multistart` | `*MultiStart` | `*Configuration` | `DefaulConfigurationtMultiStart()` |
| Simulated Annealing | `metaheuristc/sa` | `*SimulatedAnnealing` | `*Configuration` | `DefaultConfigurationSA()` |
| Genetic Algorithm | `metaheuristc/ga` | `*GA` | `*ConfigurationGA` | `DefaultConfigurationGA()` |
| BRKGA | `metaheuristc/ga` | `*BRKGA` | `*ConfigurationBRKGA` | `DefaultConfigurationBRKGA()` |
| Iterated Local Search | `metaheuristc/ils` | `*ILS` | `*Configuration` | `DefaultConfigurationILS()` |
| Variable Neighborhood Search | `metaheuristc/vns` | `*VNS` | `*Configuration` | `DefaultConfigurationVNS()` |
| Large Neighborhood Search | `metaheuristc/lns` | `*LNS` | `*Configuration` | `lns.DefaultConfigurationVNS()` |

Two exported function names do not match what they build, and are the current names in the source (not typos to fix on your side):

- `lns.DefaultConfigurationVNS()` returns an `lns.Configuration` (for LNS), not a VNS configuration.
- `multistart.DefaulConfigurationtMultiStart()` — the name itself has a typo (`Defaul...tMultiStart`).

GA and BRKGA share the `ga` subpackage but are independent types with independent `Configuration` structs; there is no shared "GA family" Configuration.

## Constructors

Every metaheuristic exposes up to three constructor shapes. All of them take `env definition.Environment`, a `*random.Generator`, a `*solution.Pool`, and a `logger.Logger`, in that order interleaved with the search/configuration arguments.

1. `CreateDefaultX(env, rg, solutionPool, logger)` — default `Configuration` and a built-in default search.
2. `CreateX(env, configuration, searchType, rg, solutionPool, logger)` — your `Configuration`, and a `search.Type` the constructor turns into a `search.Local` via `search.Create`.
3. `CreateXComplete(env, configuration, searchLocal, rg, solutionPool, logger)` — your `Configuration` and an already-built `search.Local` (e.g. a custom `search.CreateRVND` neighbourhood).

Verified signatures per subpackage:

### `multistart`

```go
func CreateDefaultMultiStart(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *MultiStart
func CreateMultiStart(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *MultiStart
```

There is no `CreateMultiStartComplete`. `CreateDefaultMultiStart` builds its search with `search.CreateDefault`.

### `sa`

```go
func CreateDefaultSA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *SimulatedAnnealing
func CreateSA(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *SimulatedAnnealing
func CreateSAComplete(env definition.Environment, config *Configuration, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) definition.Solver
```

`CreateSAComplete` returns `definition.Solver`, not `*SimulatedAnnealing` — unlike every other `...Complete` constructor in this package, which returns its concrete type. `CreateDefaultSA` builds its search with `search.CreateDefault`.

### `ga` (GA)

```go
func CreateDefaultGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *GA
func CreateGA(env definition.Environment, configuration *ConfigurationGA, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *GA
func CreateGAComplete(env definition.Environment, config *ConfigurationGA, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) *GA
```

`CreateDefaultGA` builds its search with `search.CreateMirrorLocalSearch(env)` — Mirror only, not `search.CreateDefault`'s four-neighbourhood RVND.

### `ga` (BRKGA)

```go
func CreateDefaultBRKGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *BRKGA
func CreateBRKGA(env definition.Environment, configuration *ConfigurationBRKGA, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *BRKGA
func CreateBRKGAComplete(env definition.Environment, config *ConfigurationBRKGA, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) *BRKGA
```

Same as GA: `CreateDefaultBRKGA` also uses `search.CreateMirrorLocalSearch(env)`, not `search.CreateDefault`.

### `ils`

```go
func CreateDefaultILS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *ILS
func CreateILS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *ILS
func CreateILSComplete(env definition.Environment, config *Configuration, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) *ILS
```

`CreateDefaultILS` builds its search with `search.CreateDefault`.

### `vns`

```go
func CreateDefaultVNS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *VNS
func CreateVNS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *VNS
func CreateVNSComplete(env definition.Environment, configuration *Configuration, se search.Local, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *VNS
```

`CreateDefaultVNS` builds its search with `search.CreateDefault`.

### `lns`

```go
func CreateDefaultLNS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *LNS
func CreateLNS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *LNS
func CreateLNSComplete(env definition.Environment, configuration *Configuration, se search.Local, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *LNS
```

`CreateDefaultLNS` builds its search with `search.CreateDefault`, and gets its default `Configuration` from `lns.DefaultConfigurationVNS()` (see naming note above).

## Configuration fields

Defaults come from `metaheuristc/constants/defaults.go`.

### `multistart.Configuration`

| Field | Type | Default |
|---|---|---|
| `MaxIterations` | `int` | `math.MaxInt` (unlimited) |

### `sa.Configuration`

| Field | Type | Default |
|---|---|---|
| `MaxIterations` | `int` | `math.MaxInt` (unlimited) |
| `Alpha` | `float64` | `0.0005` |
| `TemperatureInitial` | `float64` | `10000000.0` |
| `TemperatureGoal` | `float64` | `0.0000000000000001` |
| `TemperatureReheat` | `float64` | `5.0` |
| `ShakeMin` | `float64` | `0.1` |
| `ShakeMax` | `float64` | `0.3` |
| `QtdReheat` | `uint8` | `0` |
| `Iterations` | `int` | `1000` |

### `ga.ConfigurationGA`

| Field | Type | Default |
|---|---|---|
| `PopulationSize` | `int` | `200` |
| `CrossoverAlpha` | `float64` | `0.95` |
| `MutationAlpha` | `float64` | `0.005` |
| `MaxGenerations` | `int` | `math.MaxInt` (unlimited) |
| `MaxGenerationNoImprovement` | `int` | `100` |

### `ga.ConfigurationBRKGA`

| Field | Type | Default |
|---|---|---|
| `PopulationSize` | `int` | `200` |
| `EliteRatio` | `float64` | `0.2` |
| `MutantRatio` | `float64` | `0.05` |
| `CrossoverAlpha` | `float64` | `0.95` |
| `MutationAlpha` | `float64` | `0.005` |
| `MaxGenerations` | `int` | `math.MaxInt` (unlimited) |
| `MaxGenerationNoImprovement` | `int` | `100` |

### `ils.Configuration`

| Field | Type | Default |
|---|---|---|
| `MaxIterations` | `int` | `math.MaxInt` (unlimited) |
| `ShakeMin` | `float64` | `0.01` |
| `ShakeMax` | `float64` | `0.05` |
| `MetropolisCriterion` | `bool` | `false` |

### `vns.Configuration`

| Field | Type | Default |
|---|---|---|
| `MaxIterations` | `int` | `math.MaxInt` (unlimited) |
| `Rate` | `float64` | `0.2` |

### `lns.Configuration`

| Field | Type | Default |
|---|---|---|
| `MaxIterations` | `int` | `math.MaxInt` (unlimited) |
| `BetaMin` | `float64` | `0.01` |
| `BetaMax` | `float64` | `0.05` |

## Local search (`search` subpackage)

Local searches are user-selectable, both in the `Create*` constructors above (`search.Type` argument) and in the solver YAML consumed by `configuration` (see [`../configuration/README.md`](../configuration/README.md)).

### `search.Type`

```go
type Type = int

const (
    Swap Type = iota
    Mirror
    Farey
    Nelder
    RVND
)
```

| Type | Perturbs |
|---|---|
| `Swap` | Pairs of key positions within the `(start, end)` ranges from `Environment.SwapSearch()`. |
| `Mirror` | Each key value, one at a time, to `1 - value`. |
| `Farey` | Each key value, one at a time, to values drawn from sub-intervals of a fixed Farey sequence. |
| `Nelder` | The solution by blending it with other solutions drawn from the shared `solution.Pool`. |
| `RVND` | Nothing itself — it cycles through a configurable list of the searches above. |

Name/type conversion:

```go
func GetSearchType(name string) Type   // case-insensitive: "SWAP", "MIRROR", "FAREY", "NELDER", "RVND"; unknown -> -1
func GetSearchString(se Type) string   // Swap, Mirror, Farey, "Nelder-Mead", RVND; unknown -> ""
```

`GetSearchType` and `GetSearchString` are not exact inverses: the type `Nelder` round-trips through the string `"Nelder-Mead"`, not `"Nelder"`, but `GetSearchType` still accepts the input string `"Nelder"` (case-insensitive) to produce it.

### `search.Local` interface

```go
type Local interface {
    SetRG(rg *random.Generator)
    Search(ctx context.Context, rko *metaheuristc.RandomKeyValue)
    String() string
}
```

### Constructors

```go
func Create(typeSearch Type, environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local
func CreateDefault(environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local
func CreateSwapLocalSearch(environment definition.Environment) Local
func CreateMirrorLocalSearch(environment definition.Environment) Local
func CreateFareyLocalSearch(environment definition.Environment, rg *random.Generator) Local
func CreateNelderMeadLocalSearch(environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local
func CreateRVND(environment definition.Environment, s *solution.Pool, rg *random.Generator, neighbourhood []Type) Local
```

`Create` dispatches on `typeSearch`: `Swap`, `Mirror`, `Farey`, and `Nelder` each build their own single search; `RVND` builds `CreateRVND` over `[Swap, Mirror, Farey]` (three neighbourhoods, no `Nelder`); any other value (including the `-1` that `GetSearchType` returns for an unrecognized name) falls back to `CreateDefault`.

`CreateDefault` builds `CreateRVND` over `[Swap, Mirror, Farey, Nelder]` — all four neighbourhoods, including `Nelder`.

This means `search.Create(search.RVND, ...)` and `search.CreateDefault(...)` are **not** the same RVND: the first omits `Nelder`, the second includes it.

## Solve and time limits

Every solver here implements `definition.Solver`. `Solve(ctx context.Context) definition.Result` runs until `ctx` is cancelled or its deadline passes — there is no solver-local time-limit field. The time limit itself is set at the level above this package; see [`../configuration/README.md`](../configuration/README.md).

## See Also

- [Project README](../README.md)
- [`configuration` package](../configuration/README.md)
- [`definition` package](../definition/README.md)
- [`logger` package](../logger/README.md)

This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
