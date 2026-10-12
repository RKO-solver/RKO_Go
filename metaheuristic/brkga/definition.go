// Package brkga implements the Biased Random-Key Genetic Algorithm (BRKGA).
package brkga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "BRKGA"

// Config holds the BRKGA parameters.
type Config struct {
	PopulationSize             int     `yaml:"PopulationSize"`
	EliteRatio                 float64 `yaml:"EliteRatio"`
	MutantRatio                float64 `yaml:"MutantRatio"`
	CrossoverAlpha             float64 `yaml:"CrossoverAlpha"`
	MutationAlpha              float64 `yaml:"MutationAlpha"`
	MaxGenerationNoImprovement int     `yaml:"MaxGenerationNoImprovement"`
}

// BRKGA is a Biased Random-Key Genetic Algorithm solver.
type BRKGA struct {
	env           definition.Environment
	configuration *Config
	logger        logger.SolverLogger
	search        search.Local
	rg            *random.Generator
	solutionPool  *solution.Pool
}

func (brkga *BRKGA) Print() {
	search.PrintSolver(name, brkga.search)
}
