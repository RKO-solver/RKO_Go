# `config` Package

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
Parameters:
  GA:
    PopulationSize: 100
```

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/RKO-solver/rko-go/config"
)

func main() {
	cfgs, err := config.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// MyEnvironment implements definition.Environment
	solver, rkoLog, err := cfgs[0].NewSolver(MyEnvironment{})
	if err != nil {
		log.Fatal(err)
	}
	defer rkoLog.Close()

	result := solver.Solve(context.Background())
	fmt.Printf("Best solution: %+v\n", result)
}
```

`ReadFile` returns one `Config` per YAML document, so one file can hold several configurations (separate them with `---`). `Solve` returns the best solution found, passed through `Environment.Decode`.

## Loading

| Function | Use |
|---|---|
| `ReadFile(path) ([]Config, error)` | Reads every document in the file. Always strict. A file with no document, or only comments, returns an error wrapping `ErrNoConfig`. |
| `NewDecoder(r io.Reader)` | Reads YAML documents from any reader. Lenient until you call `Strict(true)`. |
| `(*Decoder).Strict(bool)` | Switches strict mode on or off. |
| `(*Decoder).Decode(*Config) error` | Decodes the next document, one per call, and returns `io.EOF` after the last. |

```go
dec := config.NewDecoder(strings.NewReader(src))
dec.Strict(true)

for {
	var cfg config.Config
	err := dec.Decode(&cfg)
	if err == io.EOF {
		break
	}
	if err != nil {
		log.Fatal(err) // do not use cfg
	}
	// use cfg
}
```

`Decode` overwrites the whole `Config`. After an error, do not use `cfg`. A document that is just `---` or `{}` gives the defaults.

### Strict and lenient

| | Strict | Lenient (default for `Decoder`) |
|---|---|---|
| Unknown or misplaced key | Error wrapping `ErrUnknownField`, with the line number | Ignored |
| Unknown `LogLevel` or `LogType` | Error wrapping `ErrUnknownValue` | Default (`INFO`, `PRINT`) |
| Unknown metaheuristic or search name | Error wrapping `ErrUnknownValue` | Dropped |
| `GRASP`, `VLNS`, `ALNS`, `IPR` (not implemented) | Error wrapping `ErrNotImplemented` | Dropped |
| Wrong type (`MaxIterations: abc`), YAML syntax error | Error | Error |
| Missing key | Default | Default |

Strict mode reports every problem at once as a joined error. Test for a kind of problem with `errors.Is`:

```go
cfgs, err := config.ReadFile("config.yaml")
switch {
case errors.Is(err, config.ErrNoConfig):
	// empty file
case errors.Is(err, config.ErrUnknownField):
	// misspelled or misplaced key; err lists each one with its line
case errors.Is(err, config.ErrUnknownValue):
	// bad LogLevel, LogType, metaheuristic or search name
case errors.Is(err, config.ErrNotImplemented):
	// GRASP, VLNS, ALNS or IPR
}
```

| Error | Meaning |
|---|---|
| `ErrNoConfig` | The file has no configuration document. |
| `ErrUnknownField` | A key is unknown or in the wrong place (strict only). |
| `ErrUnknownValue` | A value is not recognised. `NewSolver` also returns it for an invalid `LogLevel`, `LogType` or metaheuristic. |
| `ErrNotImplemented` | A recognised metaheuristic that is not implemented yet (`GRASP`, `VLNS`, `ALNS`, `IPR`). |

## Stopping

`Solve(ctx)` stops at the time limit, when every metaheuristic has run `MaxIterations` iterations, or when `ctx` is done, whichever comes first.

| `TimeLimitSeconds` | `MaxIterations` | The solver stops |
|---|---|---|
| set | not set | at the time limit |
| not set | set | when every metaheuristic has run `MaxIterations` iterations |
| set | set | at whichever comes first |
| not set | not set | only when `ctx` is cancelled |

Each metaheuristic counts its own iterations of its outer loop: a generation for GA and BRKGA, a temperature step for SA (each runs `Iterations` moves). The [`metaheuristc` package](../metaheuristc/README.md#metaheuristics) lists what one iteration is for each.

There is no default limit. With neither set, `NewSolver` prints this line to stdout:

```text
Warning: neither TimeLimitSeconds nor MaxIterations is set; the solver will run until its context is cancelled
```

## YAML reference

Keys are PascalCase and case-sensitive; values are case-insensitive. Missing keys keep their defaults.

### Top-level keys

| Key | Type | Default | Meaning |
|---|---|---|---|
| `LogLevel` | string | `INFO` | `SILENT`, `INFO` or `VERBOSE`. |
| `LogType` | string | `PRINT` | `PRINT` (plain lines, `logger/printer`) or `DASHBOARD` (live dashboard, `logger/dashboard`). |
| `TimeLimitSeconds` | int | `0` | Time limit of the whole run, in seconds. `0` or negative means no limit. |
| `MaxIterations` | int | `0` | Iteration limit of each metaheuristic. `0` or negative means no limit. |
| `SeedPool` | map | see below | `MaxSize` and `InitialSize` of the seed solution pool. |
| `Metaheuristics` | list of strings | empty | The metaheuristics to run in parallel. See below. |
| `Parameters` | map | all defaults | Tuning parameters per metaheuristic. See below. |

### `SeedPool`

The pool is shared by the metaheuristics and starts with random solutions.

| Key | Type | Default | Meaning |
|---|---|---|---|
| `MaxSize` | int | `10` | Maximum pool size. `0` or missing gives `10`; a negative value means no cap; otherwise it is raised to at least `InitialSize`. |
| `InitialSize` | int | `10` | Number of random solutions the pool starts with. `0`, negative or missing gives `10`. |

The same rules apply to a `Config` built in Go.

### `Metaheuristics` entries

Each entry is `NAME` or `NAME=search1,search2,...`, for example `"VNS=Farey,Nelder"`.

- Names: `MultiStart`, `SA`, `GA`, `BRKGA`, `ILS`, `VNS`, `LNS`.
- Searches: `Swap`, `Mirror`, `Farey`, `Nelder`. They are described in the [`metaheuristc` package](../metaheuristc/README.md#local-searches).
- Names and searches are case-insensitive, and surrounding spaces are trimmed (`" ils = swap , MIRROR "` works).
- `RVND` is ignored in search lists, since listing several searches already gives an RVND.
- `GRASP`, `VLNS`, `ALNS` and `IPR` are recognised but not implemented. See [Strict and lenient](#strict-and-lenient).
- Listing a name twice runs two independent instances.

The search list sets the entry's local search:

| Search list | Local search |
|---|---|
| none (`"GA"`, `"GA="`) | RVND over `Swap`, `Mirror`, `Farey`, `Nelder` |
| one search (`"GA=Mirror"`) | that search alone |
| two or more (`"GA=Swap,Mirror"`) | RVND over the listed searches |

### `Parameters` section

Every section and key is optional, and must sit under `Parameters:`. Field names match the Go configuration structs of the [`metaheuristc` package](../metaheuristc/README.md#configuration-structs).

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
| `QtdReheat` | int (`0`–`255`) | `0` | Number of reheats; the n-th sets the temperature to `TemperatureInitial / 2^n`. Values outside `0`–`255` are a type error. |
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

## Building a solver

```go
func (cfg *Config) NewSolver(env definition.Environment, opts ...Option) (*rko.Solver, logger.Recorder, error)
```

- Returns an error, before creating anything, for an invalid `LogLevel` or `LogType`, an invalid or unimplemented metaheuristic `Type`, or a nil section inside a non-nil `Parameters`. All problems are joined; test them with `errors.Is`.
- A nil `Parameters` means `DefaultParameters()`.
- `NewSolver` never modifies `cfg`, and each metaheuristic gets its own copy of its parameters. One `Config` can be reused for several solvers and shared across goroutines.
- The logger is named from `env.ProblemName()`. Neither `NewSolver` nor `Solve` closes it: `defer` its `Close`, which deletes the run's history files, once you have read the results. To read or save them, see the [`logger` package](../logger/README.md).
- `WithSeed(seed uint64)` seeds the random generator. Without it, each `NewSolver` call uses a randomly seeded generator.

```go
solver, rkoLog, err := cfg.NewSolver(env, config.WithSeed(42))
```

## Building a `Config` in Go

```go
params := config.DefaultParameters()
params.GA.PopulationSize = 100

cfg := config.Config{
	LogLevel:      logger.LevelInfo,
	LogType:       logger.TypePrint,
	TimeLimit:     30 * time.Second,
	MaxIterations: 1000,
	Metaheuristics: []config.Metaheuristic{
		{Type: rko.ILS, Searches: []search.Type{search.Swap, search.Mirror}},
		{Type: rko.GA}, // no searches: the default RVND
	},
	Parameters: params,
}

solver, rkoLog, err := cfg.NewSolver(env)
```

| Type | Fields |
|---|---|
| `Config` | `LogLevel`, `LogType`, `TimeLimit` (`0` or negative: no limit), `MaxIterations` (`0` or negative: no limit), `SeedPool`, `Metaheuristics`, `Parameters` |
| `SeedPool` | `MaxSize` (`0`: `solution.DefaultSize`; negative: no cap), `InitialSize` (`0` or negative: `solution.DefaultSize`) |
| `Metaheuristic` | `Type` (`rko.MetaHeuristic`), `Searches` (`[]search.Type`, same rules as the YAML table above) |
| `Parameters` | `BRKGA`, `GA`, `ILS`, `SA`, `VNS`, `LNS`: pointers to the structs of the [`metaheuristc` package](../metaheuristc/README.md#configuration-structs) |

A `Config` you build yourself skips the YAML defaults:

- The zero `LogType` is `logger.TypePrint` and the zero `LogLevel` is `logger.LevelInfo`, the same as the YAML defaults.
- `LogType` must be `logger.TypePrint` or `logger.TypeDashboard`, and `LogLevel` between `logger.LevelSilent` and `logger.LevelVerbose`.
- `Parameters` is used as is, without the greater-than-zero rule. Start from `DefaultParameters()`: a zero field stays zero, and a nil section is an error.

## See also

- [Project README](../README.md)
- [`definition` package](../definition/README.md)
- [`metaheuristc` package](../metaheuristc/README.md)
- [`logger` package](../logger/README.md)
