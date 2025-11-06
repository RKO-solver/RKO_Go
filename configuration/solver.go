package configuration

import (
	"fmt"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/channel"
	"github.com/RKO-solver/rko-go/logger/stdout"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
	"github.com/RKO-solver/rko-go/random"
)

func CreateSolver(problemName string, env definition.Environment, solverConfig *SolverConfiguration, mhConfig *MetaheuristicsConfiguration) (*rko.Solver, logger.Logger) {
	var log logger.Logger
	solvers := make([]definition.Solver, 0)

	switch solverConfig.LoggerType {
	case logger.CHANNEL:
		log = channel.NewLoggerLevel(problemName, solverConfig.LoggerLevel)
	case logger.PRINT:
		log = stdout.NewLogger(problemName, solverConfig.LoggerLevel)
	}

	rg := random.GetGlobalInstance()
	solutionPool := solution.NewDefaultPool(env, rg, log)

	for _, sol := range solverConfig.Solvers {
		var solver definition.Solver
		se := getSearch(sol.Search, env, solutionPool, rg)

		switch sol.MetaHeuristic {
		case rko.ILS:
			conf := mhConfig.ILS

			solver = ils.CreateILSComplete(env, conf, se, rg, solutionPool, log)
		case rko.VNS:
			conf := mhConfig.VNS

			solver = vns.CreateVNSComplete(env, conf, se, rg, solutionPool, log)
		case rko.MULTISTART:
			solver = multistart.CreateDefaultMultiStart(env, rg, solutionPool, log)
		case rko.SA:
			conf := mhConfig.SA

			solver = sa.CreateSAComplete(env, conf, se, rg, solutionPool, log)
		case rko.GA:
			conf := mhConfig.GA

			solver = ga.CreateGAComplete(env, conf, se, rg, solutionPool, log)
		case rko.BRKGA:
			conf := mhConfig.BRKGA

			solver = ga.CreateBRKGAComplete(env, conf, se, rg, solutionPool, log)
		case rko.LNS:
			conf := mhConfig.LNS

			solver = lns.CreateLNSComplete(env, conf, se, rg, solutionPool, log)
		default:
			fmt.Printf("%s not implemented yet\n", rko.GetMetaHeuristicString(sol.MetaHeuristic))
			continue
		}

		if solver != nil {
			solvers = append(solvers, solver)
		}
	}

	return rko.CreateFullSolver(log, rg, env, solutionPool, solvers), log
}

func CreateSolverDefaultConfig(problemName string, env definition.Environment, solverConfig *SolverConfiguration) (*rko.Solver, logger.Logger) {
	return CreateSolver(problemName, env, solverConfig, DefaultConfiguration())
}
