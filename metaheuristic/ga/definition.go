// Package ga implements the Genetic Algorithm (GA).
package ga

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const name = "GA"

// Config holds the GA parameters.
type Config struct {
	PopulationSize             metaheuristic.List[int]     `yaml:"PopulationSize"`
	CrossoverAlpha             metaheuristic.List[float64] `yaml:"CrossoverAlpha"`
	MutationAlpha              metaheuristic.List[float64] `yaml:"MutationAlpha"`
	MaxGenerationNoImprovement metaheuristic.List[int]     `yaml:"MaxGenerationNoImprovement"`
}

// GA is a Genetic Algorithm metaheuristic.
type GA struct {
	env                        problem.Environment
	configuration              *Config
	logger                     logger.MetaheuristicLogger
	search                     metaheuristic.LocalSearch
	rg                         *random.Generator
	solutionPool               *solution.Pool
	tuner                      tuning.Tuner
	states                     int
	populationSize             tuning.Param[int]
	crossoverAlpha             tuning.Param[float64]
	mutationAlpha              tuning.Param[float64]
	maxGenerationNoImprovement tuning.Param[int]
}

func (ga *GA) Print() {
	metaheuristic.Print(name, ga.search)
}
