package multistart

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/constants"
	"github.com/lucasmends/rko-go/metaheuristc/search"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

func CreateDefaultMultiStart(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *MultiStart {
	configuration := &Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
	}

	local := search.CreateDefault(env, rg)

	return &MultiStart{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateMultiStart(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *MultiStart {
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
