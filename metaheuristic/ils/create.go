package ils

import (
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const (
	defaultShakeMin = 0.01
	defaultShakeMax = 0.05
)

// DefaultConfig returns the default ILS parameters.
func DefaultConfig() *Config {
	return &Config{
		ShakeMin:            metaheuristic.List[float64]{defaultShakeMin},
		ShakeMax:            metaheuristic.List[float64]{defaultShakeMax},
		MetropolisCriterion: new(false),
	}
}

// Option configures New.
type Option func(*ILS)

// WithConfig sets the parameters; nil keeps DefaultConfig.
func WithConfig(cfg *Config) Option { return func(s *ILS) { s.configuration = cfg } }

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local metaheuristic.LocalSearch) Option { return func(s *ILS) { s.search = local } }

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *ILS) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates an ILS; without options it uses DefaultConfig and search.NewDefault.
func New(env problem.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *ILS {
	s := &ILS{
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
	if c := s.configuration; len(c.ShakeMin) > 1 || len(c.ShakeMax) > 1 {
		space := tuning.NewSpace()
		s.shakeMin = tuning.Add(space, c.ShakeMin)
		s.shakeMax = tuning.Add(space, c.ShakeMax)
		s.states = space.States()
		s.tuner = tuning.NewQLearning(space, rg)
	}
	return s
}
