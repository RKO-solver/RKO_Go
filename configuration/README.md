# `configuration` Package

Build a working RKO-Go solver from YAML files (or Go structs) instead of wiring metaheuristics by hand. This is the fastest way to get a solver running once you have a [`definition.Environment`](../definition/definition.go) implementation.

Everything here is verified against the source in this directory: `definition.go`, `create.go`, `opts.go`, `solver.go`, `helper.go`, `definition_test.go`.

## At a glance

Two independent YAML files feed this package:

| File | Loaded by | Contains |
|---|---|---|
| Metaheuristics file (e.g. `config.yaml`) | `CreateYamlMHConfiguration` | Per-metaheuristic tuning parameters (`GA`, `SA`, `VNS`, ...) |
| Solver file (e.g. `solver.yaml`) | `CreateYamlSolverConfiguration` | Logger settings, the run time limit, and which metaheuristics to run |

**Watch the casing.** The metaheuristics file uses PascalCase keys (`MaxIterations`, `PopulationSize`). The solver file uses lowerCamelCase keys (`logLevel`, `timeLimitSeconds`). Mixing them up is silently ignored by the YAML unmarshaler — no error, just a field that stays at its default.

## Data structures

```go
type MetaheuristicsConfiguration struct {
    MultiStart *multistart.Configuration
    BRKGA      *ga.ConfigurationBRKGA
    GA         *ga.ConfigurationGA
    ILS        *ils.Configuration
    SA         *sa.Configuration
    VNS        *vns.Configuration
    LNS        *lns.Configuration
}

type SolverConfiguration struct {
    LoggerLevel      logger.Level
    LoggerType       logger.LogType
    TimeLimitSeconds int64 // seconds; 0 means no time limit
    Solvers          []Solver
}

type Solver struct {
    MetaHeuristic rko.MetaHeuristic
    Search        []search.Type
}
```

There is no exported constructor for `Solver` or for building a `SolverConfiguration` by hand in a documented way — build both through `CreateYamlSolverConfiguration`, and build `MetaheuristicsConfiguration` through `DefaultConfiguration` or `CreateYamlMHConfiguration`.

## Metaheuristics YAML (`CreateYamlMHConfiguration`)

Keys are PascalCase and match the Go struct field names exactly. Every section is optional — omit a metaheuristic entirely to use all its defaults, or set only the keys you want to override.

**Merge rule:** a value from the file only overrides the default when it is greater than zero. Writing `0` (or leaving a numeric field out) silently falls back to the default; it is not treated as "explicitly zero". The one exception is `ILS.MetropolisCriterion`, a boolean: setting it to `true` in the file turns it on, but setting it to `false` (or omitting it) has no effect — there is no way to force it off once a Go-level default enables it (the current default is `false`, so this only matters if that default ever changes).

### `MultiStart`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | unlimited (`math.MaxInt`) | Iteration cap. Effectively unused — see "MultiStart is special" below. |

### `GA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `PopulationSize` | int | 200 | Number of individuals per generation |
| `CrossoverAlpha` | float64 | 0.95 | Crossover bias parameter |
| `MutationAlpha` | float64 | 0.005 | Mutation bias parameter |
| `MaxGenerations` | int | unlimited (`math.MaxInt`) | Generation cap |
| `MaxGenerationNoImprovement` | int | 100 | Stop after this many generations with no improvement |

### `BRKGA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `PopulationSize` | int | 200 | Number of individuals per generation |
| `EliteRatio` | float64 | 0.2 | Fraction of the population treated as elite |
| `MutantRatio` | float64 | 0.05 | Fraction of each new generation that is random mutants |
| `CrossoverAlpha` | float64 | 0.95 | Crossover bias parameter |
| `MutationAlpha` | float64 | 0.005 | Mutation bias parameter |
| `MaxGenerations` | int | unlimited (`math.MaxInt`) | Generation cap |
| `MaxGenerationNoImprovement` | int | 100 | Stop after this many generations with no improvement |

### `ILS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | unlimited (`math.MaxInt`) | Iteration cap |
| `ShakeMin` | float64 | 0.01 | Minimum perturbation ("shake") intensity |
| `ShakeMax` | float64 | 0.05 | Maximum perturbation ("shake") intensity |
| `MetropolisCriterion` | bool | `false` | Accept worsening moves under a Metropolis criterion; only `true` has effect (see merge rule above) |

### `SA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | unlimited (`math.MaxInt`) | Iteration cap |
| `Alpha` | float64 | 0.0005 | Cooling rate applied per iteration |
| `TemperatureInitial` | float64 | 10000000.0 | Starting temperature |
| `TemperatureGoal` | float64 | 1e-16 | Temperature at which the run is considered "cold" |
| `TemperatureReheat` | float64 | 5.0 | Temperature restored on reheat |
| `ShakeMin` | float64 | 0.1 | Minimum perturbation intensity |
| `ShakeMax` | float64 | 0.3 | Maximum perturbation intensity |
| `QtdReheat` | uint8 | 0 | Number of reheats allowed |
| `Iterations` | int | 1000 | Iterations per temperature step |

### `VNS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | unlimited (`math.MaxInt`) | Iteration cap |
| `Rate` | float64 | 0.2 | Neighborhood growth rate |

### `LNS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxIterations` | int | unlimited (`math.MaxInt`) | Iteration cap |
| `BetaMin` | float64 | 0.01 | Minimum destruction/repair fraction |
| `BetaMax` | float64 | 0.05 | Maximum destruction/repair fraction |

### Example `config.yaml`

```yaml
GA:
  PopulationSize: 100
  CrossoverAlpha: 0.7
  MutationAlpha: 0.01
  MaxGenerations: 500
  MaxGenerationNoImprovement: 50

BRKGA:
  PopulationSize: 100
  EliteRatio: 0.1
  MutantRatio: 0.1
  CrossoverAlpha: 0.7
  MutationAlpha: 0.01
  MaxGenerations: 500

ILS:
  MaxIterations: 1000
  ShakeMin: 0.02
  ShakeMax: 0.1

SA:
  MaxIterations: 2000
  TemperatureInitial: 5000
  TemperatureGoal: 0.001

VNS:
  MaxIterations: 1000
  Rate: 0.3

LNS:
  MaxIterations: 1000
```

There is no global time limit in this file — that field was removed. Time limits belong exclusively to the solver YAML below.

## Solver YAML (`CreateYamlSolverConfiguration`)

Keys are lowerCamelCase.

| Key | Type | Default | Meaning |
|---|---|---|---|
| `logLevel` | string | `INFO` | One of `SILENT`, `INFO`, `VERBOSE` (case-insensitive). Any other/unknown value also resolves to `INFO`. |
| `logType` | string | `PRINT` | One of `CHANNEL`, `PRINT` (case-insensitive). Any other/unknown value also resolves to `PRINT`. |
| `timeLimitSeconds` | int | none | Run time limit, in seconds. Omitted or `0` means the solver runs with no time limit. |
| `metaheuristics` | list of strings | none | Which metaheuristics to run, in `NAME[=search1,search2,...]` form (see below). |

### Example `solver.yaml`

```yaml
logLevel: "INFO"
logType: "CHANNEL"
timeLimitSeconds: 60
metaheuristics:
  - "GA=Swap,Mirror"
  - "SA"
  - "VNS=Farey,Nelder"
  - "ILS=Swap"
  - "BRKGA"
  - "MultiStart"
  - "LNS"
```

Following this example verbatim gives a solver bounded to 60 seconds. Drop `timeLimitSeconds` (or set it to `0`) only if you deliberately want an unbounded run.

### The `metaheuristics` entry format

Each entry is `NAME[=search1,search2,...]`. Both the metaheuristic name and the search names are matched case-insensitively.

Recognized names: `MultiStart`, `SA`, `GA`, `VNS`, `ILS`, `BRKGA`, `LNS`. (`GRASP`, `VLNS`, `ALNS`, and `IPR` are also recognized as names but have no implementation yet; listing one prints a message to stdout and that entry is skipped.)

Recognized searches: `Swap`, `Mirror`, `Farey`, `Nelder`, `RVND`.

Trip hazards, all verified in `definition.go` / `helper.go`:

- **A typo is silent.** An unrecognized metaheuristic name or search name is dropped without any warning or error. `"G A"` or `"Ga "` will simply not run.
- **Duplicates run twice.** Listing the same metaheuristic entry more than once (e.g. `"GA"` twice, or with different searches) starts two independent solver instances that both write to the shared solution pool.
- **`RVND` in a list is a no-op filter, not a real choice.** If you write `"ILS=RVND"`, the `RVND` token is filtered out while parsing that entry. Since no valid search names remain, ILS falls back to its default search — which already is RVND over all four local searches (`Swap`, `Mirror`, `Farey`, `Nelder`). So `"ILS=RVND"` and plain `"ILS"` behave identically. Use `RVND` only as an implicit default (i.e., don't list any `=search` part) rather than writing it explicitly.
- **`MultiStart` ignores everything you configure for it.** Its own `MaxIterations` setting and any `=search` list you attach to it in the `metaheuristics` entry are both ignored; it always uses `multistart.CreateDefaultMultiStart`, which picks its own default search internally.

## Go API

### `DefaultConfiguration() *MetaheuristicsConfiguration`

Returns a `MetaheuristicsConfiguration` populated with the Go-level defaults for every metaheuristic (the same defaults listed in the tables above).

### `CreateYamlMHConfiguration(filePath string) (*MetaheuristicsConfiguration, error)`

Reads the metaheuristics YAML file, starts from `DefaultConfiguration()`, and overrides fields according to the merge rule described above.

### `CreateYamlSolverConfiguration(filePath string) (*SolverConfiguration, error)`

Reads the solver YAML file and returns a `SolverConfiguration` ready to pass to `CreateSolver`.

### `CreateSolver`

```go
func CreateSolver(
    problemName string,
    env definition.Environment,
    solverConfig *SolverConfiguration,
    mhConfig *MetaheuristicsConfiguration,
) (*rko.Solver, logger.Logger)
```

Builds one solver per entry in `solverConfig.Solvers`, wired to the matching section of `mhConfig`, sharing a single solution pool and logger. Takes exactly these four arguments — there is no `random.Generator` parameter.

### `CreateSolverSeed`

```go
func CreateSolverSeed(
    problemName string,
    env definition.Environment,
    solverConfig *SolverConfiguration,
    mhConfig *MetaheuristicsConfiguration,
    seed uint64,
) (*rko.Solver, logger.Logger)
```

Same as `CreateSolver`, but runs on a random generator seeded with `seed` instead of the global instance. Use this when you need a reproducible run.

### `CreateSolverDefaultConfig`

```go
func CreateSolverDefaultConfig(
    problemName string,
    env definition.Environment,
    solverConfig *SolverConfiguration,
) (*rko.Solver, logger.Logger)
```

Same as `CreateSolver`, but uses `DefaultConfiguration()` for the metaheuristics instead of a loaded `MetaheuristicsConfiguration`. Handy when you only need to tune the solver YAML (logging, time limit, metaheuristic selection) and are fine with default tuning parameters.

### `Option`

`Option` (`type Option func(*MetaheuristicsConfiguration)`) is exported, but this package has no exported `Option` constructors — the seven `withX` functions (`withMultiStart`, `withGA`, `withBRKGA`, `withSA`, `withILS`, `withVNS`, `withLNS`) that build `Option` values are all unexported and only used internally by `CreateYamlMHConfiguration`. There is no public builder API for assembling a `MetaheuristicsConfiguration` from options; use `DefaultConfiguration()` or `CreateYamlMHConfiguration()` instead.

## Usage examples

### Fully YAML-driven

```go
package main

import (
    "fmt"
    "log"

    "github.com/RKO-solver/rko-go/configuration"
)

func main() {
    mhConfig, err := configuration.CreateYamlMHConfiguration("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    solverConfig, err := configuration.CreateYamlSolverConfiguration("solver.yaml")
    if err != nil {
        log.Fatal(err)
    }

    env := MyEnvironment{} // implements definition.Environment

    solver, _ := configuration.CreateSolver("MyProblem", env, solverConfig, mhConfig)

    result := solver.Solve()
    fmt.Printf("Best solution: %+v\n", result)
}
```

### Solver YAML only, default tuning parameters

```go
package main

import (
    "fmt"
    "log"

    "github.com/RKO-solver/rko-go/configuration"
)

func main() {
    solverConfig, err := configuration.CreateYamlSolverConfiguration("solver.yaml")
    if err != nil {
        log.Fatal(err)
    }

    env := MyEnvironment{} // implements definition.Environment

    solver, _ := configuration.CreateSolverDefaultConfig("MyProblem", env, solverConfig)

    result := solver.Solve()
    fmt.Printf("Best solution: %+v\n", result)
}
```

For a reproducible run, replace the `CreateSolver`/`CreateSolverDefaultConfig` call with:

```go
solver, _ := configuration.CreateSolverSeed("MyProblem", env, solverConfig, mhConfig, 42)
```

## Files

- `definition.go`: `MetaheuristicsConfiguration`, `SolverConfiguration`, `Solver`, `Option`, and the two YAML loaders.
- `create.go`: `DefaultConfiguration`.
- `opts.go`: unexported `withX` option functions used by `CreateYamlMHConfiguration`.
- `solver.go`: `CreateSolver`, `CreateSolverSeed`, `CreateSolverDefaultConfig`.
- `helper.go`: unexported search-selection logic used when building each solver.
- `definition_test.go`: unit tests covering the YAML merge behavior for each metaheuristic.

## See also

- [Project README](../README.md)
- [`logger` package](../logger/README.md)
- [`definition` package](../definition/README.md)
- [Metaheuristic defaults](../metaheuristc/constants/defaults.go)

This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
