package vns

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func DefaultConfigurationVNS() *Configuration {
	return &Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
		Rate:             constants.DefaultRate,
	}
}

func CreateDefaultVNS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *VNS {

	return &VNS{
		env:           env,
		configuration: DefaultConfigurationVNS(),
		logger:        logger.GetLogger(name),
		solutionPool:  solutionPool,
		RG:            rg,
	}
}

func CreateVNS(env definition.Environment, configuration *Configuration, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *VNS {
	return &VNS{
		env:           env,
		configuration: configuration,
		logger:        logger.GetLogger(name),
		solutionPool:  solutionPool,
		RG:            rg,
	}
}
