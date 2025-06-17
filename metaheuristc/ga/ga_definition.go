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
	TimeLimitSeconds           float64
	PopulationSize             int
	ChildrenRatio              float64
	CrossoverAlpha             float64
	MutationAlpha              float64
	MaxGenerations             int
	MaxGenerationNoImprovement int
}

type GA struct {
	env           definition.Environment
	configuration *ConfigurationGA
	logger        *logger.Log
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
