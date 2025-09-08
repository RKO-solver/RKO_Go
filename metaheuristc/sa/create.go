package sa

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func DefaultConfigurationSA() *Configuration {
	return &Configuration{
		MaxIterations:      constants.DefaultMaxIterations,
		TimeLimitSeconds:   constants.DefaultTimeLimitSeconds,
		Alpha:              constants.DefaultAlphaSimulationAnnealing,
		ChangeImpact:       constants.DefaultImpact,
		TemperatureInitial: constants.DefaultTemperatureInitial,
		TemperatureGoal:    constants.DefaultTemperatureGoal,
		TemperatureReheat:  constants.DefaultReheat,
		ShakeMin:           constants.DefaultShakeMinSimulationAnnealing,
		ShakeMax:           constants.DefaultShakeMaxSimulationAnnealing,
		QtdReheat:          constants.DefaultPreheat,
		Iterations:         constants.DefaultIterationsSimulationAnnealing,
	}
}

func CreateDefaultSA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *SimulatedAnnealing {
	local := search.CreateDefault(env, rg)

	return &SimulatedAnnealing{
		env:           env,
		configuration: DefaultConfigurationSA(),
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateSA(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *SimulatedAnnealing {
	local := search.Create(searchType, env, rg)
	return &SimulatedAnnealing{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
