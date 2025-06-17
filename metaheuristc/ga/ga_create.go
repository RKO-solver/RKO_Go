package ga

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/constants"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

func CreateDefaultGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *GA {
	configuration := &ConfigurationGA{
		TimeLimitSeconds:           constants.DefaultTimeLimitSeconds,
		PopulationSize:             constants.DefaultPopulationSize,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerations:             constants.DefaultMaxGenerations,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}

	local := search.CreateMirrorLocalSearch(env)

	return &GA{
		env:           env,
		configuration: configuration,
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
