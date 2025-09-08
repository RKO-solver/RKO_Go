package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func DefaultConfigurationGA() *ConfigurationGA {
	return &ConfigurationGA{
		TimeLimitSeconds:           constants.DefaultTimeLimitSeconds,
		PopulationSize:             constants.DefaultPopulationSize,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerations:             constants.DefaultMaxGenerations,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}
}

func CreateDefaultGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *GA {
	local := search.CreateMirrorLocalSearch(env)

	return &GA{
		env:           env,
		configuration: DefaultConfigurationGA(),
		logger:        logger.GetLogger(nameGA),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateGA(env definition.Environment, configuration *ConfigurationGA, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *GA {
	local := search.Create(searchType, env, rg)
	return &GA{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(nameGA),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
