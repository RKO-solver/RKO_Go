package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/dashboard"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristc/brkga"
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
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

	params := cfg.Parameters
	if params == nil {
		params = DefaultParameters()
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
		se := getSearch(m.Searches, env, pool, srg)
		s, ok := newMetaheuristic(m.Type, params, env, se, srg, pool, log)
		if !ok {
			_ = log.Close()
			return nil, nil, metaheuristicError(i, m.Type)
		}
		solvers = append(solvers, s)
	}

	return rko.CreateFullSolver(log, env, cfg.TimeLimit, cfg.MaxIterations, pool, solvers), log, nil
}

func newMetaheuristic(t rko.MetaHeuristic, p *Parameters, env definition.Environment, se search.Local, rg *random.Generator, pool *solution.Pool, log logger.Logger) (definition.Solver, bool) {
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
	case rko.MULTISTART:
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
	}
	return errors.Join(errs...)
}
