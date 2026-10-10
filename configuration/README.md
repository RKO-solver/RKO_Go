# `configuration` Package

Builds an `rko.Solver` from a YAML file or a Go struct: the logger, the stopping limits, the metaheuristics to run in parallel, their local searches and their tuning parameters. You supply a [`definition.Environment`](../definition/README.md).

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
	"context"
	"fmt"
	"log"

	"github.com/RKO-solver/rko-go/configuration"
)

func main() {
	cfg, err := configuration.CreateSolverConfiguration("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	solver, rkoLog := configuration.CreateSolver("MyProblem", MyEnvironment{}, cfg) // MyEnvironment implements definition.Environment
	defer rkoLog.Close()

	result := solver.Solve(context.Background())
	fmt.Printf("Best solution: %+v\n", result)
}
```

`Solve` returns the best solution found, passed through `Environment.Decode`.

## Stopping

`Solve(ctx)` stops at the time limit, when every metaheuristic has run `MaxIterations` iterations, or when `ctx` is done, whichever comes first.

| `TimeLimitSeconds` | `MaxIterations` | The solver stops |
|---|---|---|
| set | not set | at the time limit |
| not set | set | when every metaheuristic has run `MaxIterations` iterations |
| set | set | at whichever comes first |
| not set | not set | only when `ctx` is cancelled |

Each metaheuristic counts its own iterations of its outer loop: a generation for GA and BRKGA, a temperature step for SA (each runs `Iterations` moves). The [`metaheuristc` package](../metaheuristc/README.md#metaheuristics) lists what one iteration is for each.

There is no default limit. With neither set, `CreateSolver` prints this line to stdout:

```text
Warning: neither TimeLimitSeconds nor MaxIterations is set; the solver will run until its context is cancelled
```

## YAML reference

Keys are PascalCase and case-sensitive. Unknown or misspelled keys are ignored; missing keys keep their defaults.

### Top-level keys

| Key | Type | Default | Meaning |
|---|---|---|---|
| `LogLevel` | string | `INFO` | `SILENT`, `INFO` or `VERBOSE`, case-insensitive. Unknown values give `INFO`. |
| `LogType` | string | `PRINT` | `PRINT` (plain lines, `logger/stdout`) or `CHANNEL` (live dashboard, `logger/channel`), case-insensitive. Unknown values give `PRINT`. |
| `TimeLimitSeconds` | int | `0` | Time limit of the whole run, in seconds. `0` or negative means no limit. |
| `MaxIterations` | int | `0` | Iteration limit of each metaheuristic. `0` or negative means no limit. |
| `Metaheuristics` | list of strings | empty | The metaheuristics to run in parallel. See below. |
| `Configuration` | map | all defaults | Tuning parameters per metaheuristic. See below. |

### `Metaheuristics` entries

Each entry is `NAME` or `NAME=search1,search2,...`, for example `"VNS=Farey,Nelder"`.

- Names: `MultiStart`, `SA`, `GA`, `BRKGA`, `ILS`, `VNS`, `LNS`.
- Searches: `Swap`, `Mirror`, `Farey`, `Nelder`. They are described in the [`metaheuristc` package](../metaheuristc/README.md#local-searches).
- Names and searches are case-insensitive, and surrounding spaces are trimmed (`" ils = swap , MIRROR "` works).
- Unknown names and searches are dropped silently. `RVND` is dropped from search lists, since listing several searches already gives an RVND.
- `GRASP`, `VLNS`, `ALNS` and `IPR` are recognised but not implemented: `CreateSolver` prints `<full name> not implemented yet` and skips them.
- Listing a name twice runs two independent instances.

The search list sets the entry's local search:

| Search list | Local search |
|---|---|
| none (`"GA"`, `"GA="` or only invalid names) | RVND over `Swap`, `Mirror`, `Farey`, `Nelder` |
| one search (`"GA=Mirror"`) | that search alone |
| two or more (`"GA=Swap,Mirror"`) | RVND over the listed searches |

### `Configuration` section

Every section and key is optional, and must sit under `Configuration:` (the same keys at the top level are ignored). Field names match the Go configuration structs of the [`metaheuristc` package](../metaheuristc/README.md#configuration-structs).

A value replaces its default only when it is greater than `0`; `0` or a negative value keeps the default. `ILS.MetropolisCriterion` can only be switched on: `true` sets it, `false` keeps the default (`false`).

A **shake** applies `int(NumKeys × r) + 1` random moves to a solution, with `r` drawn from the shake range. MultiStart has no parameters.

#### `GA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `PopulationSize` | int | `200` | Individuals per generation. |
| `CrossoverAlpha` | float | `0.95` | Probability that a pair of parents is crossed. |
| `MutationAlpha` | float | `0.005` | Per-key mutation probability during a crossover. |
| `MaxGenerationNoImprovement` | int | `100` | The population restarts after more than this many generations without improvement. |

#### `BRKGA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `PopulationSize` | int | `200` | Individuals per generation. |
| `EliteRatio` | float | `0.2` | Fraction of the population kept as elite. `int(EliteRatio × PopulationSize)` must be at least `1`, or the run panics. |
| `MutantRatio` | float | `0.05` | Fraction of the population replaced by random individuals each generation. |
| `CrossoverAlpha` | float | `0.95` | Per-key probability that a child inherits the elite parent's key. |
| `MutationAlpha` | float | `0.005` | Per-key probability that a child gets a random key. |
| `MaxGenerationNoImprovement` | int | `100` | The population restarts after more than this many generations without improvement. |

#### `ILS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `ShakeMin` | float | `0.01` | Lower bound of the shake range. The range grows while the search stagnates. |
| `ShakeMax` | float | `0.05` | Upper bound of the shake range. |
| `MetropolisCriterion` | bool | `false` | Also continue from non-improving solutions, with a probability that falls to 0 at the deadline. Works only when the run has a deadline (`TimeLimitSeconds`, or a `ctx` with a deadline). |

#### `SA`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `Alpha` | float | `0.0005` | After each temperature step, the temperature is multiplied by `1 - Alpha` (unless it reheats). |
| `TemperatureInitial` | float | `10000000.0` | Starting temperature. |
| `TemperatureGoal` | float | `1e-16` | SA stops when the temperature falls to this value. |
| `TemperatureReheat` | float | `5.0` | Below this temperature, SA reheats if it has reheats left. |
| `ShakeMin` | float | `0.1` | Lower bound of the shake range. |
| `ShakeMax` | float | `0.3` | Upper bound of the shake range. |
| `QtdReheat` | int (`0`–`255`) | `0` | Number of reheats; the n-th sets the temperature to `TemperatureInitial / 2^n`. Values outside `0`–`255` fail to parse. |
| `Iterations` | int | `1000` | Moves per temperature step. |

#### `VNS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `Rate` | float | `0.2` | Neighbourhood `k` (1 to 3) shakes with `r` drawn from `[k × Rate, (k+1) × Rate)`. |

#### `LNS`

| Key | Type | Default | Meaning |
|---|---|---|---|
| `BetaMin` | float | `0.01` | Each iteration rebuilds a number of keys drawn from `[int(BetaMin × NumKeys), int(BetaMax × NumKeys))`. |
| `BetaMax` | float | `0.05` | See `BetaMin`. `int(BetaMax × NumKeys)` must be greater than `int(BetaMin × NumKeys)`, or the run panics; with the defaults, LNS needs at least 20 keys. |

## Go API

| Function | Returns |
|---|---|
| `CreateSolverConfiguration(filePath string)` | `(*SolverConfiguration, error)` read from a YAML file. The error wraps `failed to read config file` or `failed to unmarshal yaml`. |
| `CreateSolver(problemName string, env definition.Environment, solverConfig *SolverConfiguration)` | `(*rko.Solver, logger.Logger)`, using the global random generator. |
| `CreateSolverSeed(problemName string, env definition.Environment, solverConfig *SolverConfiguration, seed uint64)` | The same, with a random generator seeded with `seed`. |
| `DefaultConfiguration()` | `*MetaheuristicsConfiguration` with every metaheuristic at its defaults. |

`problemName` names the logger. The returned `logger.Logger` is the one the solver writes to: read the run's data with `GetReportData` and `GetSolutionData`, or save it to CSV (see the [`logger` package](../logger/README.md)). Neither `CreateSolver` nor `Solve` closes it: `defer` its `Close` to delete the history files once you have read the results. All metaheuristics share one solution pool and one logger; each gets its own random generator, split from the main one. Because they exchange solutions through the pool while running, a seeded run with several metaheuristics or a time limit can still vary.

### Types

```go
type SolverConfiguration struct {
    LoggerLevel   logger.Level
    LoggerType    logger.LogType
    TimeLimit     time.Duration // 0 or negative means no limit
    MaxIterations int           // per metaheuristic; 0 or negative means no limit
    Solvers       []Solver
    Configuration *MetaheuristicsConfiguration
}

type Solver struct {
    MetaHeuristic rko.MetaHeuristic
    Search        []search.Type
}

type MetaheuristicsConfiguration struct {
    BRKGA *ga.ConfigurationBRKGA
    GA    *ga.ConfigurationGA
    ILS   *ils.Configuration
    SA    *sa.Configuration
    VNS   *vns.Configuration
    LNS   *lns.Configuration
}
```

The package also declares `type Option func(*MetaheuristicsConfiguration)`, used internally by `CreateSolverConfiguration`.

### Building a configuration in Go

```go
mhConfig := configuration.DefaultConfiguration()
mhConfig.GA.PopulationSize = 100

cfg := &configuration.SolverConfiguration{
	LoggerLevel:   logger.INFO,
	LoggerType:    logger.PRINT,
	TimeLimit:     30 * time.Second,
	MaxIterations: 1000,
	Solvers: []configuration.Solver{
		{MetaHeuristic: rko.ILS, Search: []search.Type{search.Swap, search.Mirror}},
		{MetaHeuristic: rko.GA}, // no search: the default RVND
	},
	Configuration: mhConfig,
}

solver, rkoLog := configuration.CreateSolver("MyProblem", MyEnvironment{}, cfg)
defer rkoLog.Close()
```

A Go-built configuration skips the YAML parsing, so:

- Zero values are `SILENT` for `LoggerLevel` and `CHANNEL` for `LoggerType` (the YAML defaults are `INFO` and `PRINT`).
- `LoggerType` must be `logger.PRINT` or `logger.CHANNEL`; any other value leaves the logger nil and panics.
- `Configuration` is used as is, without the greater-than-zero rule. Start from `DefaultConfiguration()`: a zero field stays zero, and a nil section panics when that metaheuristic runs. A nil `Configuration` is replaced by `DefaultConfiguration()`, in your struct.
- `Search` follows the same table as YAML. A single `search.RVND` gives the default RVND; in a longer list it is skipped.

## See also

- [Project README](../README.md)
- [`definition` package](../definition/README.md)
- [`metaheuristc` package](../metaheuristc/README.md)
- [`logger` package](../logger/README.md)
