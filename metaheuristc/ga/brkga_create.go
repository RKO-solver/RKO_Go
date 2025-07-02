package ga

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/constants"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

func CreateDefaultBRKGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *BRKGA {
	configuration := &ConfigurationBRKGA{
		TimeLimitSeconds:           constants.DefaultTimeLimitSeconds,
		PopulationSize:             constants.DefaultPopulationSize,
		EliteRatio:                 constants.DefaultEliteRatio,
		MutantRation:               constants.DefaultMutantRation,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerations:             constants.DefaultMaxGenerations,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}

	local := search.CreateMirrorLocalSearch(env)

	return &BRKGA{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(nameBRKGA),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateBRKGA(env definition.Environment, configuration *ConfigurationBRKGA, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *BRKGA {
	local := search.Create(searchType, env, rg)
	return &BRKGA{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(nameBRKGA),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
