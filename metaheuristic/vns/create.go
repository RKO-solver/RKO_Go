package vns

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const defaultRate = 0.2

// DefaultConfig returns the default VNS parameters.
func DefaultConfig() *Config {
	return &Config{
		Rate: metaheuristic.List[float64]{defaultRate},
	}
}

// Option configures New.
type Option func(*VNS)

// WithConfig sets the parameters; nil keeps DefaultConfig.
func WithConfig(cfg *Config) Option { return func(s *VNS) { s.configuration = cfg } }

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local metaheuristic.LocalSearch) Option { return func(s *VNS) { s.search = local } }

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *VNS) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates a VNS; without options it uses DefaultConfig and search.NewDefault.
func New(env problem.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *VNS {
	s := &VNS{
		env:          env,
		logger:       log.MetaheuristicLogger(name),
		solutionPool: pool,
		rg:           rg,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.configuration == nil {
		s.configuration = DefaultConfig()
	}
	if s.search == nil {
		s.search = search.NewDefault(env, pool, rg)
	}
	if c := s.configuration; len(c.Rate) > 1 {
		space := tuning.NewSpace()
		s.rate = tuning.Add(space, c.Rate)
		s.states = space.States()
		s.tuner = tuning.NewQLearning(space, rg)
	}
	return s
}
