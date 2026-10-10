// Package sa implements Simulated Annealing (SA).
package sa

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "SA"

// Configuration holds the SA parameters.
type Configuration struct {
	Alpha              float64 `yaml:"Alpha"`
	TemperatureInitial float64 `yaml:"TemperatureInitial"`
	TemperatureGoal    float64 `yaml:"TemperatureGoal"`
	TemperatureReheat  float64 `yaml:"TemperatureReheat"`
	ShakeMin           float64 `yaml:"ShakeMin"`
	ShakeMax           float64 `yaml:"ShakeMax"`
	QtdReheat          uint8   `yaml:"QtdReheat"`
	Iterations         int     `yaml:"Iterations"`
}

// SimulatedAnnealing is a Simulated Annealing solver.
type SimulatedAnnealing struct {
	env           definition.Environment
	configuration *Configuration
	logger        logger.SolverLogger
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}

func (sa *SimulatedAnnealing) Print() {
	search.PrintSolver(name, sa.search)
}
