package ils

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/constants"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

func CreateDefaultILS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *ILS {
	configuration := &Configuration{
		MaxIterations:       constants.DefaultMaxIterations,
		TimeLimitSeconds:    constants.DefaultTimeLimitSeconds,
		ShakeMin:            constants.DefaultShakeMin,
		ShakeMax:            constants.DefaultShakeMax,
		MetropolisCriterion: false,
	}

	local := search.CreateDefault(env, rg)

	return &ILS{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateILS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *ILS {
	local := search.Create(searchType, env, rg)
	return &ILS{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
