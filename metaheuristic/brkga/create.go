package brkga

import (
	"fmt"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/internal/genetic"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

const (
	defaultEliteRatio  = 0.2
	defaultMutantRatio = 0.05
)

// DefaultConfig returns the default BRKGA parameters.
func DefaultConfig() *Config {
	return &Config{
		PopulationSize:             metaheuristic.List[int]{genetic.DefaultPopulationSize},
		EliteRatio:                 metaheuristic.List[float64]{defaultEliteRatio},
		MutantRatio:                metaheuristic.List[float64]{defaultMutantRatio},
		CrossoverAlpha:             metaheuristic.List[float64]{genetic.DefaultCrossoverAlpha},
		MutationAlpha:              metaheuristic.List[float64]{genetic.DefaultMutationAlpha},
		MaxGenerationNoImprovement: metaheuristic.List[int]{genetic.DefaultMaxGenerationNoImprovement},
	}
}

// Option configures New.
type Option func(*BRKGA)

// WithConfig sets the parameters; nil keeps DefaultConfig.
func WithConfig(cfg *Config) Option { return func(s *BRKGA) { s.configuration = cfg } }

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local metaheuristic.LocalSearch) Option { return func(s *BRKGA) { s.search = local } }

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *BRKGA) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates a BRKGA; without options it uses DefaultConfig and search.NewDefault.
func New(env problem.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *BRKGA {
	s := &BRKGA{
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
	for _, p := range s.configuration.PopulationSize {
		for _, e := range s.configuration.EliteRatio {
			if size := int(e * float64(p)); size < 1 || size >= p {
				panic(fmt.Sprintf("brkga: EliteRatio %v with PopulationSize %d gives elite size %d; need 1 to %d", e, p, size, p-1))
			}
		}
	}
	if c := s.configuration; len(c.PopulationSize) > 1 || len(c.EliteRatio) > 1 || len(c.MutantRatio) > 1 || len(c.CrossoverAlpha) > 1 || len(c.MutationAlpha) > 1 || len(c.MaxGenerationNoImprovement) > 1 {
		space := tuning.NewSpace()
		s.populationSize = tuning.Add(space, c.PopulationSize)
		s.eliteRatio = tuning.Add(space, c.EliteRatio)
		s.mutantRatio = tuning.Add(space, c.MutantRatio)
		s.crossoverAlpha = tuning.Add(space, c.CrossoverAlpha)
		s.mutationAlpha = tuning.Add(space, c.MutationAlpha)
		s.maxGenerationNoImprovement = tuning.Add(space, c.MaxGenerationNoImprovement)
		s.states = space.States()
		s.tuner = tuning.NewQLearning(space, rg)
	}
	return s
}
