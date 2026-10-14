package sa

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
	defaultAlpha              = 0.0005
	defaultTemperatureInitial = 10000000.0
	defaultTemperatureGoal    = 0.0000000000000001
	defaultTemperatureReheat  = 5.0
	defaultShakeMin           = 0.1
	defaultShakeMax           = 0.3
	defaultQtdReheat          = 0
	defaultIterations         = 1000
	defaultRestart            = true
)

// DefaultConfig returns the default SA parameters.
func DefaultConfig() *Config {
	return &Config{
		Alpha:              metaheuristic.List[float64]{defaultAlpha},
		TemperatureInitial: metaheuristic.List[float64]{defaultTemperatureInitial},
		TemperatureGoal:    metaheuristic.List[float64]{defaultTemperatureGoal},
		TemperatureReheat:  metaheuristic.List[float64]{defaultTemperatureReheat},
		ShakeMin:           metaheuristic.List[float64]{defaultShakeMin},
		ShakeMax:           metaheuristic.List[float64]{defaultShakeMax},
		QtdReheat:          metaheuristic.List[uint8]{defaultQtdReheat},
		Iterations:         metaheuristic.List[int]{defaultIterations},
		Restart:            new(defaultRestart),
	}
}

// Option configures New.
type Option func(*SimulatedAnnealing)

// WithConfig sets the parameters; nil keeps DefaultConfig.
func WithConfig(cfg *Config) Option { return func(s *SimulatedAnnealing) { s.configuration = cfg } }

// WithSearch sets the local search; nil keeps the default.
func WithSearch(local metaheuristic.LocalSearch) Option {
	return func(s *SimulatedAnnealing) { s.search = local }
}

// WithSearchType sets the local search to search.New(t, ...).
func WithSearchType(t search.Type) Option {
	return func(s *SimulatedAnnealing) { s.search = search.New(t, s.env, s.solutionPool, s.rg) }
}

// New creates an SA; without options it uses DefaultConfig and search.NewDefault.
func New(env problem.Environment, rg *random.Generator, pool *solution.Pool, log logger.Logger, opts ...Option) *SimulatedAnnealing {
	s := &SimulatedAnnealing{
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
	if c := s.configuration; len(c.Alpha) > 1 || len(c.TemperatureInitial) > 1 || len(c.TemperatureGoal) > 1 || len(c.TemperatureReheat) > 1 || len(c.ShakeMin) > 1 || len(c.ShakeMax) > 1 || len(c.QtdReheat) > 1 || len(c.Iterations) > 1 {
		space := tuning.NewSpace()
		s.alpha = tuning.Add(space, c.Alpha)
		s.temperatureInitial = tuning.Add(space, c.TemperatureInitial)
		s.temperatureGoal = tuning.Add(space, c.TemperatureGoal)
		s.temperatureReheat = tuning.Add(space, c.TemperatureReheat)
		s.shakeMin = tuning.Add(space, c.ShakeMin)
		s.shakeMax = tuning.Add(space, c.ShakeMax)
		s.qtdReheat = tuning.Add(space, c.QtdReheat)
		s.iterations = tuning.Add(space, c.Iterations)
		s.states = space.States()
		s.tuner = tuning.NewQLearning(space, rg)
	}
	return s
}
