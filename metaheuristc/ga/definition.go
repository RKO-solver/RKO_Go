// Package ga implements the Genetic Algorithm (GA).
package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "GA"

// ConfigurationGA holds the GA parameters.
type ConfigurationGA struct {
	PopulationSize             int     `yaml:"PopulationSize"`
	CrossoverAlpha             float64 `yaml:"CrossoverAlpha"`
	MutationAlpha              float64 `yaml:"MutationAlpha"`
	MaxGenerationNoImprovement int     `yaml:"MaxGenerationNoImprovement"`
}

// GA is a Genetic Algorithm solver.
type GA struct {
	env           definition.Environment
	configuration *ConfigurationGA
	logger        logger.SolverLogger
	search        search.Local
	rg            *random.Generator
	solutionPool  *solution.Pool
}

func (ga *GA) Print() {
	search.PrintSolver(name, ga.search)
}
