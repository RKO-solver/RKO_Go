package ga

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/constants"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultConfig returns the default GA parameters.
func DefaultConfig() *Config {
	return &Config{
		PopulationSize:             constants.DefaultPopulationSize,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}
}

// Option configures New.
type Option func(*GA)

// WithConfig sets the parameters; nil keeps DefaultConfig.
func WithConfig(cfg *Config) Option { return func(s *GA) { s.configuration = cfg } }

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local search.Local) Option { return func(s *GA) { s.search = local } }

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *GA) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates a GA; without options it uses DefaultConfig and search.NewDefault.
func New(env definition.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *GA {
	s := &GA{
		env:          env,
		logger:       log.SolverLogger(name),
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
	return s
}
