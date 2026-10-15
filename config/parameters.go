package config

import (
	"github.com/RKO-solver/rko-go/metaheuristic/brkga"
	"github.com/RKO-solver/rko-go/metaheuristic/ga"
	"github.com/RKO-solver/rko-go/metaheuristic/ils"
	"github.com/RKO-solver/rko-go/metaheuristic/lns"
	"github.com/RKO-solver/rko-go/metaheuristic/sa"
	"github.com/RKO-solver/rko-go/metaheuristic/vns"
)

// Parameters holds the tuning parameters of each metaheuristic.
type Parameters struct {
	BRKGA *brkga.Config
	GA    *ga.Config
	ILS   *ils.Config
	SA    *sa.Config
	VNS   *vns.Config
	LNS   *lns.Config
}

// DefaultParameters returns every metaheuristic's default parameters.
func DefaultParameters() *Parameters {
	return &Parameters{
		BRKGA: brkga.DefaultConfig(),
		GA:    ga.DefaultConfig(),
		ILS:   ils.DefaultConfig(),
		SA:    sa.DefaultConfig(),
		VNS:   vns.DefaultConfig(),
		LNS:   lns.DefaultConfig(),
	}
}

type parameterOption func(*Parameters)

func withBRKGA(cfg *brkga.Config) parameterOption {
	return func(c *Parameters) {
		if cfg == nil {
			return
		}

		if len(cfg.PopulationSize) > 0 {
			c.BRKGA.PopulationSize = cfg.PopulationSize
		}
		if len(cfg.EliteRatio) > 0 {
			c.BRKGA.EliteRatio = cfg.EliteRatio
		}
		if len(cfg.MutantRatio) > 0 {
			c.BRKGA.MutantRatio = cfg.MutantRatio
		}
		if len(cfg.CrossoverAlpha) > 0 {
			c.BRKGA.CrossoverAlpha = cfg.CrossoverAlpha
		}
		if len(cfg.MutationAlpha) > 0 {
			c.BRKGA.MutationAlpha = cfg.MutationAlpha
		}
		if len(cfg.MaxGenerationNoImprovement) > 0 {
			c.BRKGA.MaxGenerationNoImprovement = cfg.MaxGenerationNoImprovement
		}
	}
}

func withGA(cfg *ga.Config) parameterOption {
	return func(c *Parameters) {
		if cfg == nil {
			return
		}

		if len(cfg.PopulationSize) > 0 {
			c.GA.PopulationSize = cfg.PopulationSize
		}
		if len(cfg.CrossoverAlpha) > 0 {
			c.GA.CrossoverAlpha = cfg.CrossoverAlpha
		}
		if len(cfg.MutationAlpha) > 0 {
			c.GA.MutationAlpha = cfg.MutationAlpha
		}
		if len(cfg.MaxGenerationNoImprovement) > 0 {
			c.GA.MaxGenerationNoImprovement = cfg.MaxGenerationNoImprovement
		}

	}
}

func withSA(cfg *sa.Config) parameterOption {
	return func(c *Parameters) {
		if cfg == nil {
			return
		}

		if len(cfg.Alpha) > 0 {
			c.SA.Alpha = cfg.Alpha
		}
		if len(cfg.TemperatureInitial) > 0 {
			c.SA.TemperatureInitial = cfg.TemperatureInitial
		}
		if len(cfg.TemperatureGoal) > 0 {
			c.SA.TemperatureGoal = cfg.TemperatureGoal
		}
		if len(cfg.TemperatureReheat) > 0 {
			c.SA.TemperatureReheat = cfg.TemperatureReheat
		}
		if len(cfg.ShakeMin) > 0 {
			c.SA.ShakeMin = cfg.ShakeMin
		}
		if len(cfg.ShakeMax) > 0 {
			c.SA.ShakeMax = cfg.ShakeMax
		}
		if len(cfg.QtdReheat) > 0 {
			c.SA.QtdReheat = cfg.QtdReheat
		}
		if len(cfg.Iterations) > 0 {
			c.SA.Iterations = cfg.Iterations
		}
		if cfg.Restart != nil {
			c.SA.Restart = cfg.Restart
		}
	}
}

func withILS(cfg *ils.Config) parameterOption {
	return func(c *Parameters) {
		if cfg == nil {
			return
		}

		if len(cfg.ShakeMin) > 0 {
			c.ILS.ShakeMin = cfg.ShakeMin
		}
		if len(cfg.ShakeMax) > 0 {
			c.ILS.ShakeMax = cfg.ShakeMax
		}
		if cfg.MetropolisCriterion != nil {
			c.ILS.MetropolisCriterion = cfg.MetropolisCriterion
		}
	}
}

func withVNS(cfg *vns.Config) parameterOption {
	return func(c *Parameters) {
		if cfg == nil {
			return
		}

		if len(cfg.Rate) > 0 {
			c.VNS.Rate = cfg.Rate
		}
	}
}

func withLNS(cfg *lns.Config) parameterOption {
	return func(c *Parameters) {
		if cfg == nil {
			return
		}

		if len(cfg.BetaMin) > 0 {
			c.LNS.BetaMin = cfg.BetaMin
		}
		if len(cfg.BetaMax) > 0 {
			c.LNS.BetaMax = cfg.BetaMax
		}
	}
}
