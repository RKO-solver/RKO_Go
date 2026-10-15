// Package brkga implements the Biased Random-Key Genetic Algorithm (BRKGA).
package brkga

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const name = "BRKGA"

// Config holds the BRKGA parameters.
type Config struct {
	PopulationSize             metaheuristic.List[int]     `yaml:"PopulationSize"`
	EliteRatio                 metaheuristic.List[float64] `yaml:"EliteRatio"`
	MutantRatio                metaheuristic.List[float64] `yaml:"MutantRatio"`
	CrossoverAlpha             metaheuristic.List[float64] `yaml:"CrossoverAlpha"`
	MutationAlpha              metaheuristic.List[float64] `yaml:"MutationAlpha"`
	MaxGenerationNoImprovement metaheuristic.List[int]     `yaml:"MaxGenerationNoImprovement"`
}

// BRKGA is a Biased Random-Key Genetic Algorithm metaheuristic.
type BRKGA struct {
	env                        problem.Environment
	configuration              *Config
	logger                     logger.MetaheuristicLogger
	search                     metaheuristic.LocalSearch
	rg                         *random.Generator
	solutionPool               *solution.Pool
	tuner                      tuning.Tuner
	states                     int
	populationSize             tuning.Param[int]
	eliteRatio                 tuning.Param[float64]
	mutantRatio                tuning.Param[float64]
	crossoverAlpha             tuning.Param[float64]
	mutationAlpha              tuning.Param[float64]
	maxGenerationNoImprovement tuning.Param[int]
}

func (brkga *BRKGA) Print() {
	metaheuristic.Print(name, brkga.search)
}
