package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
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
	logger        logger.SolverLogger
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}

func (ga *GA) Print() {
	search.PrintSolver(nameGA, ga.search)
}
