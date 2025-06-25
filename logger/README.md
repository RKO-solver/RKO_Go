# `logger` Package

This package provides logging interfaces and implementations for tracking the progress and results of metaheuristic optimization in the RKO-Go library.

## Overview

Logging is essential for monitoring the optimization process, debugging, and saving results. The `logger` package defines a flexible interface and includes built-in loggers for different use cases.

## Main Components

### Interface

The core of this package is the `Interface`, which you can implement to create custom loggers:

```go
type Interface interface {
    Report(bestSolutionCost, localSolutionCost int, elapsed float64)
    Verbose(message string)
    Debug(message string)
    Info(message string)
    SetIdWorker(idWorker int)
    CreateLogger(method string) Interface
    Save()
    SaveFileName(fileName string)
}
```
- `Report`: Called to report progress (best/local solution cost, elapsed time).
- `Verbose`, `Debug`, `Info`: For logging messages at different verbosity levels.
- `SetIdWorker`: Assigns a worker/thread ID to the logger.
- `CreateLogger`: Creates a new logger for a specific method or context.
- `Save`, `SaveFileName`: For saving logs to a file.

### Built-in Loggers

- **Basic Logger** (`basic/implementation.go`): Prints progress and messages to the screen. Use `basic.CreateLogger()` to instantiate.
- **CSV Logger** (`csv/implementation.go`): Saves progress and results to a CSV file. Use `csv.CreateLogger()` to instantiate.

You can implement your own logger by satisfying the `Interface`.

## Usage Example

```go
import (
    "github.com/lucasmends/rko-go/logger/basic"
)

logger := basic.CreateLogger()
```

Pass your logger to the solver when creating it:

```go
solver := rko.CreateDefaultSolver(
    mh, env, rko.INFO, false, logger,
)
```

## Files
- `definition.go`: Contains the `Interface` definition.
- `basic/implementation.go`: Basic logger implementation (prints to screen).
- `csv/implementation.go`: CSV logger implementation (saves to file).

## See Also
- [Project README](../README.md)
- [definition package](../definition/README.md)

---
This package is part of the [RKO-Go](https://github.com/lucasmends/rko-go) library.
