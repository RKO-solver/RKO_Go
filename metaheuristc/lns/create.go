package lns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func DefaultConfigurationVNS() *Configuration {
	return &Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
		BetaMin:          constants.DefaultShakeMin,
		BetaMax:          constants.DefaultShakeMax,
	}
}

func CreateDefaultLNS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *LNS {
	local := search.CreateDefault(env, solutionPool, rg)

	return &LNS{
		env:           env,
		configuration: DefaultConfigurationVNS(),
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateLNS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *LNS {
	local := search.Create(searchType, env, solutionPool, rg)

	return &LNS{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateLNSComplete(env definition.Environment, configuration *Configuration, se search.Local, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *LNS {
	return &LNS{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		search:        se,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
