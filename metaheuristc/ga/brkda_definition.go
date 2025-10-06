package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const nameBRKGA = "BRKGA"

type ConfigurationBRKGA struct {
	TimeLimitSeconds           float64
	PopulationSize             int     `yaml:"PopulationSize"`
	EliteRatio                 float64 `yaml:"EliteRatio"`
	MutantRatio                float64 `yaml:"MutantRatio"`
	CrossoverAlpha             float64 `yaml:"CrossoverAlpha"`
	MutationAlpha              float64 `yaml:"MutationAlpha"`
	MaxGenerations             int     `yaml:"MaxGenerations"`
	MaxGenerationNoImprovement int     `yaml:"MaxGenerationNoImprovement"`
}

type BRKGA struct {
	env           definition.Environment
	configuration *ConfigurationBRKGA
	logger        logger.SolverLogger
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
