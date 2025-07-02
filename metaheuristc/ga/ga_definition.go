package ga

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

const nameGA = "GA"

type ConfigurationGA struct {
	TimeLimitSeconds           float64 `yaml:"TimeLimitSeconds"`
	PopulationSize             int     `yaml:"PopulationSize"`
	ChildrenRatio              float64 `yaml:"ChildrenRatio"`
	CrossoverAlpha             float64 `yaml:"CrossoverAlpha"`
	MutationAlpha              float64 `yaml:"MutationAlpha"`
	MaxGenerations             int     `yaml:"MaxGenerations"`
	MaxGenerationNoImprovement int     `yaml:"MaxGenerationNoImprovement"`
}

type GA struct {
	env           definition.Environment
	configuration *ConfigurationGA
	logger        *logger.Log
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
