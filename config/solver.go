package config

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/dashboard"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristic/brkga"
	"github.com/RKO-solver/rko-go/metaheuristic/ga"
	"github.com/RKO-solver/rko-go/metaheuristic/ils"
	"github.com/RKO-solver/rko-go/metaheuristic/lns"
	"github.com/RKO-solver/rko-go/metaheuristic/multistart"
	"github.com/RKO-solver/rko-go/metaheuristic/sa"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/vns"
	"github.com/RKO-solver/rko-go/random"
)

// NewSolver builds a solver and its logger from cfg; the caller closes the logger.
func (cfg *Config) NewSolver(env definition.Environment, opts ...Option) (*rko.Solver, logger.Recorder, error) {
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}

	var o options
	for _, opt := range opts {
		opt(&o)
	}

	params := DefaultParameters()
	if cfg.Parameters != nil {
		params = cfg.Parameters.merged()
	}

	maxSize, initialSize := seedPoolSizes(cfg.SeedPool.MaxSize, cfg.SeedPool.InitialSize)

	var log logger.Recorder
	switch cfg.LogType {
	case logger.TypeDashboard:
		log = dashboard.New(env.ProblemName(), &dashboard.Options{Level: cfg.LogLevel})
	default:
		log = printer.New(os.Stdout, env.ProblemName(), &printer.Options{Level: cfg.LogLevel})
	}

	rg := o.generator()
	pool := solution.NewPoolSize(env, rg, log, maxSize, initialSize)

	solvers := make([]definition.Solver, 0, len(cfg.Metaheuristics))
	for i, m := range cfg.Metaheuristics {
		srg := rg.Split()
		senv := env
		if c, ok := env.(definition.Cloner); ok {
			senv = c.Clone()
		}
		se := getSearch(m.Searches, senv, pool, srg)
		s, ok := newMetaheuristic(m.Type, params, senv, se, srg, pool, log)
		if !ok {
			_ = log.Close()
			return nil, nil, metaheuristicError(i, m.Type)
		}
		solvers = append(solvers, s)
	}

	return rko.NewFromSolvers(env, log, pool, solvers, rko.WithTimeLimit(cfg.TimeLimit), rko.WithMaxIterations(cfg.MaxIterations)), log, nil
}

func newMetaheuristic(t rko.Metaheuristic, p *Parameters, env definition.Environment, se search.Local, rg *random.Generator, pool *solution.Pool, log logger.Logger) (definition.Solver, bool) {
	switch t {
	case rko.ILS:
		c := *p.ILS
		return ils.New(env, rg, pool, log, ils.WithConfig(&c), ils.WithSearch(se)), true
	case rko.VNS:
		c := *p.VNS
		return vns.New(env, rg, pool, log, vns.WithConfig(&c), vns.WithSearch(se)), true
	case rko.SA:
		c := *p.SA
		return sa.New(env, rg, pool, log, sa.WithConfig(&c), sa.WithSearch(se)), true
	case rko.GA:
		c := *p.GA
		return ga.New(env, rg, pool, log, ga.WithConfig(&c), ga.WithSearch(se)), true
	case rko.BRKGA:
		c := *p.BRKGA
		return brkga.New(env, rg, pool, log, brkga.WithConfig(&c), brkga.WithSearch(se)), true
	case rko.LNS:
		c := *p.LNS
		return lns.New(env, rg, pool, log, lns.WithConfig(&c), lns.WithSearch(se)), true
	case rko.MultiStart:
		return multistart.New(env, rg, pool, log, multistart.WithSearch(se)), true
	}
	return nil, false
}

func (cfg *Config) validate() error {
	var errs []error
	if cfg.LogLevel < logger.LevelSilent || cfg.LogLevel > logger.LevelVerbose {
		errs = append(errs, fmt.Errorf("config: LogLevel: %w %d", ErrUnknownValue, cfg.LogLevel))
	}
	if cfg.LogType > logger.TypeDashboard {
		errs = append(errs, fmt.Errorf("config: LogType: %w %d", ErrUnknownValue, cfg.LogType))
	}
	for i, m := range cfg.Metaheuristics {
		if err := metaheuristicError(i, m.Type); err != nil {
			errs = append(errs, err)
		}
	}
	if p := cfg.Parameters; p != nil {
		for _, s := range []struct {
			name string
			nil  bool
		}{
			{"BRKGA", p.BRKGA == nil},
			{"GA", p.GA == nil},
			{"ILS", p.ILS == nil},
			{"SA", p.SA == nil},
			{"VNS", p.VNS == nil},
			{"LNS", p.LNS == nil},
		} {
			if s.nil {
				errs = append(errs, fmt.Errorf("config: Parameters.%s is nil", s.name))
			}
		}
		if len(errs) == 0 {
			errs = append(errs, p.merged().validate()...)
		}
	}
	return errors.Join(errs...)
}

// merged returns the default parameters overridden by the non-empty lists of p.
func (p *Parameters) merged() *Parameters {
	m := DefaultParameters()
	for _, opt := range []parameterOption{
		withGA(p.GA),
		withBRKGA(p.BRKGA),
		withSA(p.SA),
		withVNS(p.VNS),
		withILS(p.ILS),
		withLNS(p.LNS),
	} {
		opt(m)
	}
	return m
}

// validate checks the lists of p, which must have every section and no empty list.
func (p *Parameters) validate() []error {
	var errs []error
	for _, c := range []struct {
		name string
		dup  bool
	}{
		{"BRKGA.PopulationSize", hasDuplicates(p.BRKGA.PopulationSize)},
		{"BRKGA.EliteRatio", hasDuplicates(p.BRKGA.EliteRatio)},
		{"BRKGA.MutantRatio", hasDuplicates(p.BRKGA.MutantRatio)},
		{"BRKGA.CrossoverAlpha", hasDuplicates(p.BRKGA.CrossoverAlpha)},
		{"BRKGA.MutationAlpha", hasDuplicates(p.BRKGA.MutationAlpha)},
		{"BRKGA.MaxGenerationNoImprovement", hasDuplicates(p.BRKGA.MaxGenerationNoImprovement)},
		{"GA.PopulationSize", hasDuplicates(p.GA.PopulationSize)},
		{"GA.CrossoverAlpha", hasDuplicates(p.GA.CrossoverAlpha)},
		{"GA.MutationAlpha", hasDuplicates(p.GA.MutationAlpha)},
		{"GA.MaxGenerationNoImprovement", hasDuplicates(p.GA.MaxGenerationNoImprovement)},
		{"ILS.ShakeMin", hasDuplicates(p.ILS.ShakeMin)},
		{"ILS.ShakeMax", hasDuplicates(p.ILS.ShakeMax)},
		{"SA.Alpha", hasDuplicates(p.SA.Alpha)},
		{"SA.TemperatureInitial", hasDuplicates(p.SA.TemperatureInitial)},
		{"SA.TemperatureReheat", hasDuplicates(p.SA.TemperatureReheat)},
		{"SA.ShakeMin", hasDuplicates(p.SA.ShakeMin)},
		{"SA.ShakeMax", hasDuplicates(p.SA.ShakeMax)},
		{"SA.QtdReheat", hasDuplicates(p.SA.QtdReheat)},
		{"SA.Iterations", hasDuplicates(p.SA.Iterations)},
		{"VNS.Rate", hasDuplicates(p.VNS.Rate)},
		{"LNS.BetaMin", hasDuplicates(p.LNS.BetaMin)},
		{"LNS.BetaMax", hasDuplicates(p.LNS.BetaMax)},
	} {
		if c.dup {
			errs = append(errs, fmt.Errorf("config: Parameters.%s: %w: repeated value", c.name, ErrInvalidValue))
		}
	}
	for _, c := range []struct {
		name     string
		min, max []float64
	}{
		{"ILS.ShakeMin", p.ILS.ShakeMin, p.ILS.ShakeMax},
		{"SA.ShakeMin", p.SA.ShakeMin, p.SA.ShakeMax},
		{"LNS.BetaMin", p.LNS.BetaMin, p.LNS.BetaMax},
	} {
		if lo, hi := slices.Max(c.min), slices.Min(c.max); lo > hi {
			errs = append(errs, fmt.Errorf("config: Parameters.%s: %w: %v is greater than the maximum %v", c.name, ErrInvalidValue, lo, hi))
		}
	}
	if len(p.SA.TemperatureGoal) > 1 {
		errs = append(errs, fmt.Errorf("config: Parameters.SA.TemperatureGoal: %w: only one value is allowed", ErrInvalidValue))
	}
	if len(p.ILS.MetropolisCriterion) > 1 {
		errs = append(errs, fmt.Errorf("config: Parameters.ILS.MetropolisCriterion: %w: only one value is allowed", ErrInvalidValue))
	}
	return errs
}

func hasDuplicates[T comparable](values []T) bool {
	for i := range values {
		if slices.Contains(values[:i], values[i]) {
			return true
		}
	}
	return false
}
