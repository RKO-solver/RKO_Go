package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultConfigurationBRKGA returns the default BRKGA parameters.
func DefaultConfigurationBRKGA() *ConfigurationBRKGA {
	return &ConfigurationBRKGA{
		PopulationSize:             constants.DefaultPopulationSize,
		EliteRatio:                 constants.DefaultEliteRatio,
		MutantRatio:                constants.DefaultMutantRatio,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}
}

// CreateDefaultBRKGA creates a BRKGA with the default configuration and local search.
func CreateDefaultBRKGA(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *BRKGA {
	local := search.CreateDefault(env, solutionPool, rg)

	return &BRKGA{
		env:           env,
		configuration: DefaultConfigurationBRKGA(),
		logger:        logger.GetLogger(nameBRKGA),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

// CreateBRKGA creates a BRKGA with the local search of searchType.
func CreateBRKGA(env definition.Environment, configuration *ConfigurationBRKGA, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *BRKGA {
	local := search.Create(searchType, env, solutionPool, rg)
	return &BRKGA{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(nameBRKGA),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

// CreateBRKGAComplete creates a BRKGA with the local search se.
func CreateBRKGAComplete(env definition.Environment, config *ConfigurationBRKGA, se search.Local, rg *random.Generator, pool *solution.Pool, l logger.Logger) *BRKGA {
	return &BRKGA{
		env:           env,
		configuration: config,
		logger:        l.GetLogger(nameBRKGA),
		search:        se,
		solutionPool:  pool,
		RG:            rg,
	}
}
