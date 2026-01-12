# `configuration` Package

This package provides utilities for configuring and instantiating metaheuristic solvers in the RKO-Go library. It supports both programmatic configuration and YAML-based configuration files.

## Overview

The `configuration` package simplifies solver setup by allowing you to:
- Define metaheuristic parameters via YAML configuration files
- Specify solver configurations and logger settings
- Manage time limits and algorithm-specific parameters for all metaheuristics
- Support multiple search strategies for different metaheuristics

## Main Components

### Data Structures

#### MetaheuristicsConfiguration
Holds configuration for all metaheuristics:

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
```

#### SolverConfiguration
Specifies solver behavior and logger settings:

```go
type SolverConfiguration struct {
    LoggerLevel logger.Level
    LoggerType  logger.LogType
    Solvers     []Solver
}
```

#### Solver
Defines which metaheuristic to use and its search strategies:

```go
type Solver struct {
    MetaHeuristic rko.MetaHeuristic
    Search        []search.Type
}
```

### Key Functions

#### `DefaultConfiguration()`
Creates a `MetaheuristicsConfiguration` with default parameters for all metaheuristics.

```go
config := configuration.DefaultConfiguration()
```

#### `CreateYamlMHConfiguration(filePath string)`
Loads metaheuristic configurations from a YAML file. Supports specifying a global `TimeLimitSeconds` parameter that applies to all metaheuristics.

```go
config, err := configuration.CreateYamlMHConfiguration("config.yaml")
if err != nil {
    log.Fatal(err)
}
```

#### `CreateYamlSolverConfiguration(filePath string)`
Loads solver configuration from a YAML file, including logger settings and metaheuristic definitions.

```go
solverConfig, err := configuration.CreateYamlSolverConfiguration("solver.yaml")
if err != nil {
    log.Fatal(err)
}
```

#### `CreateSolver()`
Factory function that creates a configured solver instance with appropriate logger and metaheuristics.

### Configuration Options

The package uses the Option pattern to allow flexible configuration modifications. Options are builder functions that modify a `MetaheuristicsConfiguration`:

```go
type Option func(*MetaheuristicsConfiguration)
```

Options include: `withMultiStart`, `withGA`, `withBRKGA`, `withSA`, `withVNS`, `withILS`, `withLNS`

## YAML Configuration Format

### Metaheuristics Configuration (`config.yaml`)

```yaml
TimeLimitSeconds: 60.0

MultiStart:
  MaxIterations: 100

GA:
  PopulationSize: 100
  EliteRatio: 0.1
  MutantRatio: 0.1
  CrossoverAlpha: 0.7
  MaxGenerations: 100

BRKGA:
  PopulationSize: 100
  EliteRatio: 0.1
  MutantRatio: 0.1
  CrossoverAlpha: 0.7
  MaxGenerations: 100

ILS:
  MaxIterations: 100
  PerturbationStrength: 0.1

SA:
  InitialTemperature: 100.0
  TemperatureCoolingRate: 0.95
  MaxIterations: 1000

VNS:
  MaxNeighborhoodSize: 5
  MaxIterations: 100

LNS:
  MaxIterations: 100
```

### Solver Configuration (`solver.yaml`)

```yaml
logLevel: "INFO"          # SILENT, INFO, VERBOSE
logType: "CHANNEL"        # CHANNEL, PRINT
metaheuristics:
  - "GA=Swap,Mirror"
  - "SA"
  - "VNS=Farey,Nelder"
  - "ILS=Swap"
  - "BRKGA"
  - "MultiStart"
  - "LNS"
```

**Metaheuristics Format:** `METAHEURISTIC_NAME[=search1,search2,...]`

**Available Search Strategies:**
- `Swap`: Swap-based neighborhood search
- `Mirror`: Mirror-based neighborhood search
- `Farey`: Farey sequence-based search
- `Nelder`: Nelder-Mead local search
- `RVND`: Random Variable Neighborhood Descent (automatic combination)

## Usage Example

### Programmatic Configuration

```go
package main

import (
    "github.com/RKO-solver/rko-go"
    "github.com/RKO-solver/rko-go/configuration"
    "github.com/RKO-solver/rko-go/logger/channel"
)

func main() {
    // Create default configuration
    config := configuration.DefaultConfiguration()
    
    // Your environment implementation
    env := MyEnvironment{}
    
    // Create and run solver
    mh := []rko.MetaHeuristic{rko.GA, rko.SA, rko.VNS}
    logger := channel.DefaultLogger("MyProblem")
    
    solver := rko.CreateDefaultSolver(mh, env, logger)
    result := solver.Solve()
}
```

### YAML-Based Configuration

```go
package main

import (
    "github.com/RKO-solver/rko-go/configuration"
)

func main() {
    // Load configurations from files
    mhConfig, err := configuration.CreateYamlMHConfiguration("config.yaml")
    if err != nil {
        panic(err)
    }
    
    solverConfig, err := configuration.CreateYamlSolverConfiguration("solver.yaml")
    if err != nil {
        panic(err)
    }
    
    // Create solver with loaded configurations
    env := MyEnvironment{}
    rg := random.GetGlobalInstance()
    
    solver, logger := configuration.CreateSolver(
        "MyProblem",
        env,
        solverConfig,
        mhConfig,
        rg,
    )
    
    result := solver.Solve()
}
```

## Files

- `definition.go`: Core data structures and YAML loading functions.
- `create.go`: Factory function for creating default configurations.
- `opts.go`: Option functions for flexible configuration.
- `solver.go`: Solver creation and instantiation logic.
- `helper.go`: Helper functions for search strategy selection.
- `definition_test.go`: Unit tests for configuration functions.

## See Also

- [Project README](../README.md)
- [logger package](../logger/README.md)
- [definition package](../definition/README.md)
- [Metaheuristics](../metaheuristc/)

---

This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.

