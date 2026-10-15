package rko

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/brkga"
	"github.com/RKO-solver/rko-go/metaheuristic/ga"
	"github.com/RKO-solver/rko-go/metaheuristic/ils"
	"github.com/RKO-solver/rko-go/metaheuristic/lns"
	"github.com/RKO-solver/rko-go/metaheuristic/multistart"
	"github.com/RKO-solver/rko-go/metaheuristic/sa"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/vns"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

var (
	// ErrNoMetaheuristics is returned when a Solver is given no metaheuristics.
	ErrNoMetaheuristics = errors.New("no metaheuristics")
	// ErrNotImplemented marks a metaheuristic.Type that is not implemented yet.
	ErrNotImplemented = errors.New("not implemented")
	// ErrUnknownMetaheuristic marks a value that is not a metaheuristic.Type.
	ErrUnknownMetaheuristic = errors.New("unknown metaheuristic")
)

// Option configures New and NewFromMetaheuristics.
type Option func(*Solver)

// WithTimeLimit sets the time limit; zero or negative means no limit.
func WithTimeLimit(d time.Duration) Option { return func(s *Solver) { s.timeLimit = d } }

// WithMaxIterations sets the iteration limit; zero or negative means no limit.
func WithMaxIterations(n int) Option { return func(s *Solver) { s.maxIterations = n } }

// WithMetaheuristics adds metaheuristics for New to run; repeated calls append.
func WithMetaheuristics(types ...metaheuristic.Type) Option {
	return func(s *Solver) { s.types = append(s.types, types...) }
}

// New creates a Solver running the WithMetaheuristics types, limited to DefaultTimeLimit.
func New(env problem.Environment, log logger.Logger, opts ...Option) (*Solver, error) {
	s := &Solver{l: log, env: env, timeLimit: DefaultTimeLimit}
	for _, opt := range opts {
		opt(s)
	}
	if len(s.types) == 0 {
		return nil, fmt.Errorf("rko: %w", ErrNoMetaheuristics)
	}
	for _, t := range s.types { // validate before any side effect on rg, env or log
		if t.String() == "" {
			return nil, fmt.Errorf("rko: %w %d", ErrUnknownMetaheuristic, t)
		}
		if !t.Implemented() {
			return nil, fmt.Errorf("rko: %s is %w", t, ErrNotImplemented)
		}
	}

	rg := random.NewRandom()
	solutionPool := solution.NewPool(env, rg, log)

	metaheuristics := make([]metaheuristic.Algorithm, 0, len(s.types))

	for _, t := range s.types {
		var mh metaheuristic.Algorithm
		srg := rg.Split()
		senv := env
		if c, ok := env.(problem.Cloner); ok {
			senv = c.Clone()
		}
		switch t {
		case metaheuristic.ILS:
			mh = ils.New(senv, srg, solutionPool, log)
		case metaheuristic.VNS:
			mh = vns.New(senv, srg, solutionPool, log)
		case metaheuristic.MultiStart:
			mh = multistart.New(senv, srg, solutionPool, log)
		case metaheuristic.SA:
			mh = sa.New(senv, srg, solutionPool, log)
		case metaheuristic.GA:
			mh = ga.New(senv, srg, solutionPool, log)
		case metaheuristic.BRKGA:
			mh = brkga.New(senv, srg, solutionPool, log)
		case metaheuristic.LNS:
			mh = lns.New(senv, srg, solutionPool, log)
		}

		if mh != nil {
			metaheuristics = append(metaheuristics, mh)
		}
	}

	s.solutionPool = solutionPool
	s.metaheuristics = metaheuristics
	s.types = nil
	return build(s, nil), nil
}

// NewFromMetaheuristics creates a Solver from already-built metaheuristics and their pool; without options it has no limit.
func NewFromMetaheuristics(env problem.Environment, log logger.Logger, pool *solution.Pool, metaheuristics []metaheuristic.Algorithm, opts ...Option) (*Solver, error) {
	if len(metaheuristics) == 0 {
		return nil, fmt.Errorf("rko: %w", ErrNoMetaheuristics)
	}
	return build(&Solver{l: log, env: env, solutionPool: pool, metaheuristics: metaheuristics}, opts), nil
}

func build(s *Solver, opts []Option) *Solver {
	for _, opt := range opts {
		opt(s)
	}
	if s.timeLimit <= 0 && s.maxIterations <= 0 {
		fmt.Fprintln(os.Stdout, NoLimitWarning)
	}
	return s
}
