package sa

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultConfigurationSA returns the default SA parameters.
func DefaultConfigurationSA() *Configuration {
	return &Configuration{
		Alpha:              constants.DefaultAlphaSimulationAnnealing,
		TemperatureInitial: constants.DefaultTemperatureInitial,
		TemperatureGoal:    constants.DefaultTemperatureGoal,
		TemperatureReheat:  constants.DefaultReheat,
		ShakeMin:           constants.DefaultShakeMinSimulationAnnealing,
		ShakeMax:           constants.DefaultShakeMaxSimulationAnnealing,
		QtdReheat:          constants.DefaultPreheat,
		Iterations:         constants.DefaultIterationsSimulationAnnealing,
	}
}

// CreateDefaultSA creates a SimulatedAnnealing with the default configuration and local search.
func CreateDefaultSA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *SimulatedAnnealing {
	local := search.CreateDefault(env, solutionPool, rg)

	return &SimulatedAnnealing{
		env:           env,
		configuration: DefaultConfigurationSA(),
		logger:        logger.SolverLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		rg:            rg,
	}
}

// CreateSA creates a SimulatedAnnealing with the local search of searchType.
func CreateSA(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *SimulatedAnnealing {
	local := search.Create(searchType, env, solutionPool, rg)
	return &SimulatedAnnealing{
		env:           env,
		configuration: configuration,
		logger:        logger.SolverLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		rg:            rg,
	}
}

// CreateSAComplete creates a SimulatedAnnealing with the local search se.
func CreateSAComplete(env definition.Environment, config *Configuration, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) definition.Solver {
	return &SimulatedAnnealing{
		env:           env,
		configuration: config,
		logger:        l.SolverLogger(name),
		search:        se,
		solutionPool:  pool,
		rg:            rg,
	}
}
