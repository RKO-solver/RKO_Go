package vns

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/metaheuristc/constants"
	"github.com/lucasmends/rko-go/metaheuristc/solution"
	"github.com/lucasmends/rko-go/random"
)

func CreateDefaultVNS(env definition.Environment, rg *random.Generator, solutionPool *solution.Pool, logger *logger.Log) *VNS {
	configuration := &Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
		Rate:             constants.DefaultRate,
	}

	return &VNS{
		env:           env,
		configuration: configuration,
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
