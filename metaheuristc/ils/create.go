package ils

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func DefaultConfigurationILS() *Configuration {
	return &Configuration{
		MaxIterations:       constants.DefaultMaxIterations,
		TimeLimitSeconds:    constants.DefaultTimeLimitSeconds,
		ShakeMin:            constants.DefaultShakeMin,
		ShakeMax:            constants.DefaultShakeMax,
		MetropolisCriterion: false,
	}
}

func CreateDefaultILS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *ILS {
	local := search.CreateDefault(env, rg)

	return &ILS{
		env:           env,
		configuration: DefaultConfigurationILS(),
		logger:        logger.GetLogger(name),
		search:        local,
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateILS(env definition.Environment, configuration *Configuration, searchType search.Type, rg *random.Generator, solutionPool *solution.Pool, logger logger.Logger) *ILS {
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
