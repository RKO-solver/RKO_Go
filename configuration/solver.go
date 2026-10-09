package configuration

import (
	"fmt"
	"io"
	"os"

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

// noLimitWarning is printed when a solver has neither a time limit nor an iteration limit.
const noLimitWarning = "Warning: neither TimeLimitSeconds nor MaxIterations is set; the solver will run until its context is cancelled"

// warnIfUnbounded writes noLimitWarning to w when neither TimeLimit nor MaxIterations is positive.
func warnIfUnbounded(w io.Writer, solverConfig *SolverConfiguration) {
	if solverConfig.TimeLimit <= 0 && solverConfig.MaxIterations <= 0 {
		fmt.Fprintln(w, noLimitWarning)
	}
}

// createSolver builds the logger, shared pool and one solver per implemented entry, warning when no limit is set; a nil Configuration is set to DefaultConfiguration().
func createSolver(problemName string, env definition.Environment, solverConfig *SolverConfiguration, rg *random.Generator) (*rko.Solver, logger.Logger) {
	warnIfUnbounded(os.Stdout, solverConfig)
	var log logger.Logger
	solvers := make([]definition.Solver, 0)
	if solverConfig.Configuration == nil {
		solverConfig.Configuration = DefaultConfiguration()
	}
	mhConfig := solverConfig.Configuration
	switch solverConfig.LoggerType {
	case logger.CHANNEL:
		log = channel.NewLoggerLevel(problemName, solverConfig.LoggerLevel)
	case logger.PRINT:
		log = stdout.NewLogger(problemName, solverConfig.LoggerLevel)
	}

	solutionPool := solution.NewDefaultPool(env, rg, log)

	for _, sol := range solverConfig.Solvers {
		var solver definition.Solver
		srg := rg.Split()
		se := getSearch(sol.Search, env, solutionPool, srg)

		switch sol.MetaHeuristic {
		case rko.ILS:
			conf := mhConfig.ILS

			solver = ils.CreateILSComplete(env, conf, se, srg, solutionPool, log)
		case rko.VNS:
			conf := mhConfig.VNS

			solver = vns.CreateVNSComplete(env, conf, se, srg, solutionPool, log)
		case rko.MULTISTART:
			solver = multistart.CreateMultiStart(env, se, srg, solutionPool, log)
		case rko.SA:
			conf := mhConfig.SA

			solver = sa.CreateSAComplete(env, conf, se, srg, solutionPool, log)
		case rko.GA:
			conf := mhConfig.GA

			solver = ga.CreateGAComplete(env, conf, se, srg, solutionPool, log)
		case rko.BRKGA:
			conf := mhConfig.BRKGA

			solver = ga.CreateBRKGAComplete(env, conf, se, srg, solutionPool, log)
		case rko.LNS:
			conf := mhConfig.LNS

			solver = lns.CreateLNSComplete(env, conf, se, srg, solutionPool, log)
		default:
			fmt.Printf("%s not implemented yet\n", rko.GetMetaHeuristicString(sol.MetaHeuristic))
			continue
		}

		if solver != nil {
			solvers = append(solvers, solver)
		}
	}

	return rko.CreateFullSolver(log, env, solverConfig.TimeLimit, solverConfig.MaxIterations, solutionPool, solvers), log
}

// CreateSolver builds a solver from solverConfig using the global random generator.
func CreateSolver(problemName string, env definition.Environment, solverConfig *SolverConfiguration) (*rko.Solver, logger.Logger) {
	return createSolver(problemName, env, solverConfig, random.GetGlobalInstance())
}

// CreateSolverSeed is like CreateSolver but uses a random generator seeded with seed.
func CreateSolverSeed(problemName string, env definition.Environment, solverConfig *SolverConfiguration, seed uint64) (*rko.Solver, logger.Logger) {
	rg := random.NewGeneratorSeed(seed)
	return createSolver(problemName, env, solverConfig, rg)
}
