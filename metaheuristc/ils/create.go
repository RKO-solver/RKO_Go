package ils

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultConfigurationILS returns the default ILS parameters.
func DefaultConfigurationILS() *Configuration {
	return &Configuration{
		ShakeMin:            constants.DefaultShakeMin,
		ShakeMax:            constants.DefaultShakeMax,
		MetropolisCriterion: false,
	}
}

// CreateDefaultILS creates an ILS with the default configuration and local search.
func CreateDefaultILS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *ILS {
	local := search.NewDefault(env, solutionPool, rg)

	return &ILS{
		env:           env,
		configuration: DefaultConfigurationILS(),
		logger:        logger.SolverLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		rg:            rg,
	}
}

// CreateILS creates an ILS with the local search of searchType.
func CreateILS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *ILS {
	local := search.New(searchType, env, solutionPool, rg)
	return &ILS{
		env:           env,
		configuration: configuration,
		logger:        logger.SolverLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		rg:            rg,
	}
}

// CreateILSComplete creates an ILS with the local search se.
func CreateILSComplete(env definition.Environment, config *Configuration, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) *ILS {
	return &ILS{
		env:           env,
		configuration: config,
		logger:        l.SolverLogger(name),
		search:        se,
		solutionPool:  pool,
		rg:            rg,
	}
}
