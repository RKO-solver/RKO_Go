package rko

import (
	"fmt"
	"os"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/brkga"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
	"github.com/RKO-solver/rko-go/random"
)

// Option configures New and NewFromSolvers.
type Option func(*Solver)

// WithTimeLimit sets the time limit; zero or negative means no limit.
func WithTimeLimit(d time.Duration) Option { return func(s *Solver) { s.timeLimit = d } }

// WithMaxIterations sets the iteration limit; zero or negative means no limit.
func WithMaxIterations(n int) Option { return func(s *Solver) { s.maxIterations = n } }

// New creates a Solver that runs the metaheuristics in mhs with default parameters; without options the time limit is constants.DefaultTimeLimit.
func New(env definition.Environment, log logger.Logger, mhs []Metaheuristic, opts ...Option) *Solver {
	rg := random.NewRandom()
	solutionPool := solution.NewPool(env, rg, log)

	solvers := make([]definition.Solver, 0, len(mhs))

	for _, m := range mhs {
		var solver definition.Solver
		srg := rg.Split()
		switch m {
		case ILS:
			solver = ils.New(env, srg, solutionPool, log)
		case VNS:
			solver = vns.New(env, srg, solutionPool, log)
		case MultiStart:
			solver = multistart.New(env, srg, solutionPool, log)
		case SA:
			solver = sa.New(env, srg, solutionPool, log)
		case GA:
			solver = ga.New(env, srg, solutionPool, log)
		case BRKGA:
			solver = brkga.New(env, srg, solutionPool, log)
		case LNS:
			solver = lns.New(env, srg, solutionPool, log)
		default:
			fmt.Printf("%s not implemented yet\n", m)
			continue
		}

		if solver != nil {
			solvers = append(solvers, solver)
		}
	}

	return build(&Solver{
		l:            log,
		env:          env,
		solutionPool: solutionPool,
		solvers:      solvers,
		timeLimit:    constants.DefaultTimeLimit,
	}, opts)
}

// NewFromSolvers creates a Solver from already-built solvers and their pool; without options it has no limit.
func NewFromSolvers(env definition.Environment, log logger.Logger, pool *solution.Pool, solvers []definition.Solver, opts ...Option) *Solver {
	return build(&Solver{l: log, env: env, solutionPool: pool, solvers: solvers}, opts)
}

// build applies opts to s and prints NoLimitWarning when s ends up with neither limit.
func build(s *Solver, opts []Option) *Solver {
	for _, opt := range opts {
		opt(s)
	}
	if s.timeLimit <= 0 && s.maxIterations <= 0 {
		fmt.Fprintln(os.Stdout, NoLimitWarning)
	}
	return s
}
