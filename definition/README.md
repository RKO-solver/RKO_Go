# `definition` Package

This package defines the core interfaces and types required to use the RKO-Go library for optimization problems.

## Overview

The `definition` package provides the abstraction layer between your problem and the metaheuristics implemented in RKO-Go. By implementing the interfaces in this package, you can use any of the metaheuristics (GA, SA, ILS, VNS, MultiStart, BRKGA, etc.) provided by the library.

## Main Components

### Environment Interface

The central interface is `Environment`, which you must implement for your specific optimization problem:

```go
type Environment interface {
    NumKeys() int
    Cost(r RandomKey) int
    Decode(r RandomKey) any
}
```
- `NumKeys() int`: Returns the number of random keys (variables) in your problem.
- `Cost(r RandomKey) int`: Evaluates the cost (objective function) for a given solution.
- `Decode(r RandomKey) any`: Converts a solution from random keys to your problem's representation.

### RandomKey Type

A `RandomKey` is a slice of float64 values representing a solution in the random key encoding. This encoding is used by all metaheuristics in RKO-Go for solution representation and manipulation.

## Usage

To use RKO-Go with your problem:
1. Implement the `Environment` interface in your own type.
2. Pass your implementation to the solver creation functions (e.g., `rko.CreateDefaultSolver`).

## Example

```go
package myproblem

import "github.com/RKO-solver/rko-go/definition"

type MyEnv struct{}

func (e MyEnv) NumKeys() int {
    return 10 // Number of variables
}

func (e MyEnv) Cost(r definition.RandomKey) int {
    // Compute and return the cost
}

func (e MyEnv) Decode(r definition.RandomKey) any {
    // Convert random key to your solution representation
}
```

## Files
- `definition.go`: Contains the `Environment` interface and related types.

## See Also
- [Project README](../README.md)
- [logger package](../logger/README.md)

---
This package is part of the [RKO-Go](https://github.com/RKO-solver/rko-go) library.
