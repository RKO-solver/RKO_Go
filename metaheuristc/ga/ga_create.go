package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultConfigurationGA returns the default GA parameters.
func DefaultConfigurationGA() *ConfigurationGA {
	return &ConfigurationGA{
		PopulationSize:             constants.DefaultPopulationSize,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}
}

// CreateDefaultGA creates a GA with the default configuration and local search.
func CreateDefaultGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *GA {
	local := search.NewDefault(env, solutionPool, rg)

	return &GA{
		env:           env,
		configuration: DefaultConfigurationGA(),
		logger:        logger.SolverLogger(nameGA),
		search:        local,
		solutionPool:  solutionPool,
		rg:            rg,
	}
}

// CreateGA creates a GA with the local search of searchType.
func CreateGA(env definition.Environment, configuration *ConfigurationGA, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *GA {
	local := search.New(searchType, env, solutionPool, rg)
	return &GA{
		env:           env,
		configuration: configuration,
		logger:        logger.SolverLogger(nameGA),
		search:        local,
		solutionPool:  solutionPool,
		rg:            rg,
	}
}

// CreateGAComplete creates a GA with the local search se.
func CreateGAComplete(env definition.Environment, config *ConfigurationGA, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) *GA {
	return &GA{
		env:           env,
		configuration: config,
		logger:        l.SolverLogger(nameGA),
		search:        se,
		solutionPool:  pool,
		rg:            rg,
	}
}
