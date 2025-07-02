package multistart

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

const name = "MultiStart"

type Configuration struct {
	MaxIterations    int     `yaml:"MaxIterations"`
	TimeLimitSeconds float64 `yaml:"TimeLimitSeconds"`
}

type MultiStart struct {
	env           definition.Environment
	configuration *Configuration
	logger        *logger.Log
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
