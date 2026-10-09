# `configuration` Package

Builds an `rko.Solver` from a YAML file or from Go structs: logger, time limit, the metaheuristics to run in parallel, their local searches and their tuning parameters. You supply a [`definition.Environment`](../definition/README.md).

## Quick start

`config.yaml`:

```yaml
LogLevel: INFO
LogType: PRINT
TimeLimitSeconds: 60
Metaheuristics:
  - "ILS=Swap,Mirror"
  - "GA"
Configuration:
  GA:
    PopulationSize: 100
```

```go
package main

import (
    "fmt"
    "log"

    "github.com/RKO-solver/rko-go/configuration"
)

func main() {
    cfg, err := configuration.CreateSolverConfiguration("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    solver, _ := configuration.CreateSolver("MyProblem", MyEnvironment{}, cfg) // MyEnvironment implements definition.Environment

    result := solver.Solve()
    fmt.Printf("Best solution: %+v\n", result)
}
```

`Solve` stops at the time limit; with no limit it runs until every metaheuristic stops on its own. `SolveCtx(ctx)` ignores the configured limit and stops when `ctx` is cancelled or expires. Both return the best solution decoded by `Environment.Decode`.

## YAML reference

Keys are PascalCase and case-sensitive. Unknown or misspelled keys are silently ignored; missing keys keep their defaults.

### Top-level keys

| Key | Type | Default | Meaning |
|---|---|---|---|
| `LogLevel` | string | `INFO` | `SILENT`, `INFO` or `VERBOSE`, case-insensitive. Unknown values give `INFO`. |
| `LogType` | string | `PRINT` | `CHANNEL` (live dashboard) or `PRINT` (plain lines), case-insensitive. Unknown values give `PRINT`. |
| `TimeLimitSeconds` | int | `0` | Time limit in seconds; `0` or negative means no limit. Stored as a `time.Duration` in `SolverConfiguration.TimeLimit`. |
| `Metaheuristics` | list of strings | empty | Metaheuristics to run in parallel. See the entry format below. |
| `Configuration` | map | all defaults | Per-metaheuristic tuning parameters. See below. |

### `Metaheuristics` entries

Each entry is `NAME[=search1,search2,...]`, for example `"VNS=Farey,Nelder"`.

- Names and searches are case-insensitive, and surrounding whitespace is trimmed (`" ils = swap , MIRROR "` works).
- Names: `MultiStart`, `SA`, `GA`, `VNS`, `ILS`, `BRKGA`, `LNS`.
- Searches: `Swap`, `Mirror`, `Farey`, `Nelder`.
- Unknown names and unknown searches are dropped without a warning.
- `RVND` is dropped from search lists; RVND is what you get by listing no searches or several.
- `GRASP`, `VLNS`, `ALNS` and `IPR` parse but have no implementation. `CreateSolver` prints `<full name> not implemented yet` to stdout and skips them.
- Listing a name twice runs two independent instances.

The search list decides each entry's local search:

| Search list | Local search used |
|---|---|
| none (`"GA"`, `"GA="`, or only invalid entries) | RVND over `Swap`, `Mirror`, `Farey`, `Nelder` |
| one search (`"GA=Mirror"`) | that search alone |
| two or more (`"GA=Swap,Mirror"`) | RVND over the listed searches |

### `Configuration` section

Every section and key is optional. Sections must sit under `Configuration:`; the same keys at the top level are ignored.

**Merge rule:** a numeric value overrides its default only when it is greater than `0`. `0` or a negative value keeps the default. `ILS.MetropolisCriterion` can only be switched on: `true` sets it, `false` has no effect.

`math.MaxInt` means no cap; the run is then bounded by the time limit or the context.

A shake applies `int(NumKeys × r) + 1` random perturbations. ILS and SA draw `r` from `[ShakeMin, ShakeMax)`; VNS neighbourhood `k` (1 to 3) draws it from `[k × Rate, (k+1) × Rate)`.

#### `MultiStart`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | `math.MaxInt` | Iteration cap |

#### `GA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `PopulationSize` | int | `200` | Individuals per generation |
| `CrossoverAlpha` | float64 | `0.95` | Probability that a pair of parents is crossed |
| `MutationAlpha` | float64 | `0.005` | Per-key mutation probability during crossover |
| `MaxGenerations` | int | `math.MaxInt` | Generation cap |
| `MaxGenerationNoImprovement` | int | `100` | The population is restarted once more than this many generations pass without improvement |

#### `BRKGA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `PopulationSize` | int | `200` | Individuals per generation |
| `EliteRatio` | float64 | `0.2` | Elite fraction of the population |
| `MutantRatio` | float64 | `0.05` | Fraction of the population replaced by random individuals each generation |
| `CrossoverAlpha` | float64 | `0.95` | Per-key probability that a child keeps the elite parent's key (when the key is not mutated) |
| `MutationAlpha` | float64 | `0.005` | Per-key probability of a random key in a child |
| `MaxGenerations` | int | `math.MaxInt` | Generation cap |
| `MaxGenerationNoImprovement` | int | `100` | The population is restarted once more than this many generations pass without improvement |

#### `ILS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | `math.MaxInt` | Iteration cap |
| `ShakeMin` | float64 | `0.01` | Lower bound of the shake range |
| `ShakeMax` | float64 | `0.05` | Upper bound of the shake range |
| `MetropolisCriterion` | bool | `false` | Accept non-improving neighbours as the current solution with a probability that falls to 0 at the deadline. Active only when the run context has a deadline (`Solve` with a time limit, or `SolveCtx` with a deadline). |

The shake range adapts: while more than 10 and fewer than 1000 consecutive iterations have passed without improvement, each shake grows both bounds by 5%; if either exceeds `1.0`, both reset to 1.1 times their configured values. An improvement restores the configured values.

#### `SA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `Alpha` | float64 | `0.0005` | Temperature is multiplied by `1 - Alpha` after each block of `Iterations`, unless a reheat happens |
| `TemperatureInitial` | float64 | `10000000.0` | Starting temperature |
| `TemperatureGoal` | float64 | `1e-16` | The run stops once the temperature falls to this value |
| `TemperatureReheat` | float64 | `5.0` | A reheat triggers when the temperature drops below this value |
| `ShakeMin` | float64 | `0.1` | Lower bound of the shake range |
| `ShakeMax` | float64 | `0.3` | Upper bound of the shake range |
| `QtdReheat` | uint8 | `0` | Number of reheats; the n-th sets the temperature to `TemperatureInitial / 2^n`. Values outside `0..255` fail to parse. |
| `Iterations` | int | `1000` | Iterations per temperature step |

`MaxIterations` is accepted under `SA` but ignored.

#### `VNS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | `math.MaxInt` | Iteration cap |
| `Rate` | float64 | `0.2` | Width of each neighbourhood's shake range (see above) |

#### `LNS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | `math.MaxInt` | Iteration cap |
| `BetaMin` | float64 | `0.01` | The number of keys rebuilt per iteration is drawn from `[int(BetaMin × NumKeys), int(BetaMax × NumKeys))` |
| `BetaMax` | float64 | `0.05` | See `BetaMin` |

Defaults come from [`metaheuristc/constants/defaults.go`](../metaheuristc/constants/defaults.go).

## Go API

### Types

```go
type SolverConfiguration struct {
    LoggerLevel   logger.Level
    LoggerType    logger.LogType
    TimeLimit     time.Duration // 0 or negative means no limit
    Solvers       []Solver
    Configuration *MetaheuristicsConfiguration
}

type Solver struct {
    MetaHeuristic rko.MetaHeuristic
    Search        []search.Type
}

type MetaheuristicsConfiguration struct {
    MultiStart *multistart.Configuration
    BRKGA      *ga.ConfigurationBRKGA
    GA         *ga.ConfigurationGA
    ILS        *ils.Configuration
    SA         *sa.Configuration
    VNS        *vns.Configuration
    LNS        *lns.Configuration
}

type Option func(*MetaheuristicsConfiguration)
```

No exported function builds or accepts an `Option`; `CreateSolverConfiguration` uses its internal implementations (`withGA`, `withSA`, ...) to apply the merge rule.

### Functions

| Function | Returns |
|---|---|
| `DefaultConfiguration()` | `*MetaheuristicsConfiguration` with every metaheuristic at its defaults |
| `CreateSolverConfiguration(filePath string)` | `(*SolverConfiguration, error)`; the error wraps `failed to read config file` or `failed to unmarshal yaml` |
| `CreateSolver(problemName string, env definition.Environment, solverConfig *SolverConfiguration)` | `(*rko.Solver, logger.Logger)`, using the global random generator |
| `CreateSolverSeed(problemName string, env definition.Environment, solverConfig *SolverConfiguration, seed uint64)` | Same as `CreateSolver`, with a generator seeded with `seed` |

`problemName` is passed to the logger. The returned `logger.Logger` is the one the solver writes to. All metaheuristics share one solution pool, one logger and one random generator. `TimeLimit` is passed unchanged to `rko.CreateFullSolver`.

### Building a configuration in Go

```go
cfg := &configuration.SolverConfiguration{
    LoggerLevel: logger.INFO,
    LoggerType:  logger.PRINT,
    TimeLimit:   30 * time.Second,
    Solvers: []configuration.Solver{
        {MetaHeuristic: rko.ILS, Search: []search.Type{search.Swap, search.Mirror}},
        {MetaHeuristic: rko.GA}, // no searches: default RVND
    },
    // Configuration: nil means DefaultConfiguration()
}

solver, _ := configuration.CreateSolver("MyProblem", MyEnvironment{}, cfg)
```

Go-built configurations skip YAML parsing, so:

- A single `search.RVND` gives the default RVND; in a longer list it is skipped.
- Zero values are `LoggerLevel` = `SILENT`, `LoggerType` = `CHANNEL` (the YAML default is `PRINT`), `TimeLimit` = no limit.
- `LoggerType` must be `logger.CHANNEL` or `logger.PRINT`. Any other value leaves the logger nil: `CreateSolver` panics on the first implemented metaheuristic, or `Solve`/`SolveCtx` panics if there is none.
- `Configuration` is used as is, without the merge rule. Start from `DefaultConfiguration()` and change fields: zero-valued fields stay zero, and a nil section crashes the process when that metaheuristic starts.

## Gotchas

- **`CreateSolver` and `CreateSolverSeed` mutate their input.** A nil `solverConfig.Configuration` is set to `DefaultConfiguration()` in the caller's struct.
- **Defaults differ from `rko.CreateDefaultSolver`.** There, GA and BRKGA use the `Mirror` search and the time limit is 150 s. Here, every metaheuristic without a search list uses the default RVND, and the time limit defaults to none.
- **`CreateSolverSeed` fixes only the random sequence.** All metaheuristics draw from one unsynchronized generator concurrently, so a run with two or more of them is not reproducible.
- **Some parameter values panic at run time.** LNS panics when `int(BetaMax × NumKeys) <= int(BetaMin × NumKeys)`; with the defaults, that is any problem with fewer than 20 keys. BRKGA panics when `int(EliteRatio × PopulationSize)` is `0` and `MutantRatio` is below `1` (for example `EliteRatio: 0.001` with the default `PopulationSize`).

## Files

- `definition.go`: exported types, the `Option` type and the YAML mirror structs.
- `create.go`: `DefaultConfiguration`, `CreateSolverConfiguration`, the `Metaheuristics` entry parser and the option merge.
- `opts.go`: internal `Option` implementations applying the merge rule.
- `solver.go`: `CreateSolver` and `CreateSolverSeed`.
- `helper.go`: search-list resolution.
- `create_test.go`: tests for YAML parsing, `Metaheuristics` entries, defaults and the merge rule.

## See also

- [Project README](../README.md)
- [`definition` package](../definition/README.md)
- [`logger` package](../logger/README.md)
- [`metaheuristc` package](../metaheuristc/README.md)
