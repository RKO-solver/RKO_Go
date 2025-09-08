package sa

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const name = "SA"

type Configuration struct {
	MaxIterations      int     `yaml:"MaxIterations"`
	TimeLimitSeconds   float64 `yaml:"TimeLimitSeconds"`
	Alpha              float64 `yaml:"Alpha"`
	ChangeImpact       float64 `yaml:"ChangeImpact"`
	TemperatureInitial float64 `yaml:"TemperatureInitial"`
	TemperatureGoal    float64 `yaml:"TemperatureGoal"`
	TemperatureReheat  float64 `yaml:"TemperatureReheat"`
	ShakeMin           float64 `yaml:"ShakeMin"`
	ShakeMax           float64 `yaml:"ShakeMax"`
	QtdReheat          uint8   `yaml:"QtdReheat"`
	Iterations         int     `yaml:"Iterations"`
}

type SimulatedAnnealing struct {
	env           definition.Environment
	configuration *Configuration
	logger        *logger.Log
	search        search.Local
	RG            *random.Generator
	solutionPool  *solution.Pool
}
