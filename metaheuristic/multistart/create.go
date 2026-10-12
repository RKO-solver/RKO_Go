package multistart

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
)

// Option configures New.
type Option func(*MultiStart)

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local search.Local) Option { return func(s *MultiStart) { s.search = local } }

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *MultiStart) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates a MultiStart; without options it uses search.NewDefault.
func New(env definition.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *MultiStart {
	s := &MultiStart{
		env:          env,
		logger:       log.SolverLogger(name),
		solutionPool: pool,
		rg:           rg,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.search == nil {
		s.search = search.NewDefault(env, pool, rg)
	}
	return s
}
