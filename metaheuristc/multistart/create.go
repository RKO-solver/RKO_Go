package multistart

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func DefaulConfigurationtMultiStart() *Configuration {
	return &Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
	}
}

func CreateDefaultMultiStart(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *MultiStart {
	local := search.CreateDefault(env, rg)

	return &MultiStart{
		env:           env,
		configuration: DefaulConfigurationtMultiStart(),
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateMultiStart(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *MultiStart {
	local := search.Create(searchType, env, rg)
	return &MultiStart{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
