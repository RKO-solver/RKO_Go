// Package config builds solvers from Go structs or YAML files.
package config

import (
	"errors"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
)

var (
	// ErrNoConfig marks a file with no configuration document.
	ErrNoConfig = errors.New("no configuration found")
	// ErrUnknownField marks an unknown or misplaced key in strict mode.
	ErrUnknownField = errors.New("unknown field")
	// ErrUnknownValue marks an unrecognised value.
	ErrUnknownValue = errors.New("unknown value")
	// ErrNotImplemented marks a recognised metaheuristic that is not implemented yet.
	ErrNotImplemented = errors.New("not implemented")
)

// Config describes one solver: logging, limits, seed pool and metaheuristics.
type Config struct {
	LogLevel       logger.Level  // zero value is logger.LevelInfo
	LogType        logger.Type   // zero value is logger.TypePrint
	TimeLimit      time.Duration // 0 or negative means no limit
	MaxIterations  int           // 0 or negative means no limit
	SeedPool       SeedPool
	Metaheuristics []Metaheuristic
	Parameters     *Parameters // nil means DefaultParameters()
}

// SeedPool sizes the solution pool shared by the metaheuristics.
type SeedPool struct {
	MaxSize     int // 0 means solution.DefaultSize; negative means no cap
	InitialSize int // 0 or negative means solution.DefaultSize
}

// Metaheuristic is one metaheuristic to run and its local searches.
type Metaheuristic struct {
	Type     rko.Metaheuristic
	Searches []search.Type
}
