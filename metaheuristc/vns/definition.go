package vns

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

const name = "VNS"

type Configuration struct {
	MaxIterations    int     `yaml:"MaxIterations"`
	TimeLimitSeconds float64 `yaml:"TimeLimitSeconds"`
	Rate             float64 `yaml:"Rate"`
}

type VNS struct {
	env           definition.Environment
	configuration *Configuration
	logger        *logger.Log
	RG            *random.Generator
	solutionPool  *solution.Pool
}
