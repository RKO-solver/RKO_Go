// Package sa implements Simulated Annealing (SA).
package sa

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const name = "SA"

// Config holds the SA parameters.
type Config struct {
	Alpha              metaheuristic.List[float64] `yaml:"Alpha"`
	TemperatureInitial metaheuristic.List[float64] `yaml:"TemperatureInitial"`
	TemperatureGoal    metaheuristic.List[float64] `yaml:"TemperatureGoal"`
	TemperatureReheat  metaheuristic.List[float64] `yaml:"TemperatureReheat"`
	ShakeMin           metaheuristic.List[float64] `yaml:"ShakeMin"`
	ShakeMax           metaheuristic.List[float64] `yaml:"ShakeMax"`
	QtdReheat          metaheuristic.List[uint8]   `yaml:"QtdReheat"`
	Iterations         metaheuristic.List[int]     `yaml:"Iterations"`
	Restart            *bool                       `yaml:"Restart"`
}

// SimulatedAnnealing is a Simulated Annealing metaheuristic.
type SimulatedAnnealing struct {
	env                problem.Environment
	configuration      *Config
	logger             logger.MetaheuristicLogger
	search             metaheuristic.LocalSearch
	rg                 *random.Generator
	solutionPool       *solution.Pool
	tuner              tuning.Tuner
	states             int
	alpha              tuning.Param[float64]
	temperatureInitial tuning.Param[float64]
	temperatureGoal    tuning.Param[float64]
	temperatureReheat  tuning.Param[float64]
	shakeMin           tuning.Param[float64]
	shakeMax           tuning.Param[float64]
	qtdReheat          tuning.Param[uint8]
	iterations         tuning.Param[int]
}

func (sa *SimulatedAnnealing) Print() {
	metaheuristic.Print(name, sa.search)
}
