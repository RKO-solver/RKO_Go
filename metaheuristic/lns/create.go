package lns

import (
	"fmt"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/constants"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultConfig returns the default LNS parameters.
func DefaultConfig() *Config {
	return &Config{
		BetaMin: metaheuristic.List[float64]{constants.DefaultShakeMin},
		BetaMax: metaheuristic.List[float64]{constants.DefaultShakeMax},
	}
}

// Option configures New.
type Option func(*LNS)

// WithConfig sets the parameters; nil keeps DefaultConfig.
func WithConfig(cfg *Config) Option { return func(s *LNS) { s.configuration = cfg } }

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local search.Local) Option { return func(s *LNS) { s.search = local } }

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *LNS) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates an LNS; without options it uses DefaultConfig and search.NewDefault.
func New(env definition.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *LNS {
	s := &LNS{
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
	n := env.NumKeys()
	for _, a := range s.configuration.BetaMin {
		for _, b := range s.configuration.BetaMax {
			if int(a*float64(n)) >= int(b*float64(n)) {
				panic(fmt.Sprintf("lns: BetaMin %v and BetaMax %v give an empty intensity range with %d keys", a, b, n))
			}
		}
	}
	if c := s.configuration; len(c.BetaMin) > 1 || len(c.BetaMax) > 1 {
		space := tuning.NewSpace()
		s.betaMin = tuning.Add(space, c.BetaMin)
		s.betaMax = tuning.Add(space, c.BetaMax)
		s.states = space.States()
		s.tuner = tuning.NewQLearning(space, rg)
	}
	return s
}
