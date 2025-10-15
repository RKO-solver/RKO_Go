package configuration

import (
	"math"

	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
)

func withMultiStart(cfg *multistart.Configuration, timeLimitSeconds float64) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if timeLimitSeconds > 0 {
			c.MultiStart.TimeLimitSeconds = timeLimitSeconds
		} else if timeLimitSeconds == 0 {
			c.MultiStart.TimeLimitSeconds = math.MaxFloat64
		}

		if cfg.MaxIterations > 0 {
			c.MultiStart.MaxIterations = cfg.MaxIterations
		}
	}
}

func withBRKGA(cfg *ga.ConfigurationBRKGA, timeLimitSeconds float64) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if timeLimitSeconds > 0 {
			c.BRKGA.TimeLimitSeconds = timeLimitSeconds
		} else if timeLimitSeconds == 0 {
			c.BRKGA.TimeLimitSeconds = math.MaxFloat64
		}

		if cfg.PopulationSize > 0 {
			c.BRKGA.PopulationSize = cfg.PopulationSize
		}
		if cfg.EliteRatio > 0 {
			c.BRKGA.EliteRatio = cfg.EliteRatio
		}
		if cfg.MutantRatio > 0 {
			c.BRKGA.MutantRatio = cfg.MutantRatio
		}
		if cfg.CrossoverAlpha > 0 {
			c.BRKGA.CrossoverAlpha = cfg.CrossoverAlpha
		}
		if cfg.MutationAlpha > 0 {
			c.BRKGA.MutationAlpha = cfg.MutationAlpha
		}
		if cfg.MaxGenerations > 0 {
			c.BRKGA.MaxGenerations = cfg.MaxGenerations
		}
		if cfg.MaxGenerationNoImprovement > 0 {
			c.BRKGA.MaxGenerationNoImprovement = cfg.MaxGenerationNoImprovement
		}
	}
}

func withGA(cfg *ga.ConfigurationGA, timeLimitSeconds float64) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if timeLimitSeconds > 0 {
			c.GA.TimeLimitSeconds = timeLimitSeconds
		} else if timeLimitSeconds == 0 {
			c.GA.TimeLimitSeconds = math.MaxFloat64
		}

		if cfg.PopulationSize > 0 {
			c.GA.PopulationSize = cfg.PopulationSize
		}
		if cfg.ChildrenRatio > 0 {
			c.GA.ChildrenRatio = cfg.ChildrenRatio
		}
		if cfg.CrossoverAlpha > 0 {
			c.GA.CrossoverAlpha = cfg.CrossoverAlpha
		}
		if cfg.MutationAlpha > 0 {
			c.GA.MutationAlpha = cfg.MutationAlpha
		}
		if cfg.MaxGenerations > 0 {
			c.GA.MaxGenerations = cfg.MaxGenerations
		}
		if cfg.MaxGenerationNoImprovement > 0 {
			c.GA.MaxGenerationNoImprovement = cfg.MaxGenerationNoImprovement
		}

	}
}

func withSA(cfg *sa.Configuration, timeLimitSeconds float64) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if timeLimitSeconds > 0 {
			c.SA.TimeLimitSeconds = timeLimitSeconds
		} else if timeLimitSeconds == 0 {
			c.SA.TimeLimitSeconds = math.MaxFloat64
		}

		if cfg.MaxIterations > 0 {
			c.SA.MaxIterations = cfg.MaxIterations
		}
		if cfg.Alpha > 0 {
			c.SA.Alpha = cfg.Alpha
		}
		if cfg.ChangeImpact > 0 {
			c.SA.ChangeImpact = cfg.ChangeImpact
		}
		if cfg.TemperatureInitial > 0 {
			c.SA.TemperatureInitial = cfg.TemperatureInitial
		}
		if cfg.TemperatureGoal > 0 {
			c.SA.TemperatureGoal = cfg.TemperatureGoal
		}
		if cfg.TemperatureReheat > 0 {
			c.SA.TemperatureReheat = cfg.TemperatureReheat
		}
		if cfg.ShakeMin > 0 {
			c.SA.ShakeMin = cfg.ShakeMin
		}
		if cfg.ShakeMax > 0 {
			c.SA.ShakeMax = cfg.ShakeMax
		}
		if cfg.QtdReheat > 0 {
			c.SA.QtdReheat = cfg.QtdReheat
		}
		if cfg.Iterations > 0 {
			c.SA.Iterations = cfg.Iterations
		}
	}
}

func withILS(cfg *ils.Configuration, timeLimitSeconds float64) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if timeLimitSeconds > 0 {
			c.ILS.TimeLimitSeconds = timeLimitSeconds
		} else if timeLimitSeconds == 0 {
			c.ILS.TimeLimitSeconds = math.MaxFloat64
		}

		if cfg.MaxIterations > 0 {
			c.ILS.MaxIterations = cfg.MaxIterations
		}
		if cfg.ShakeMin > 0 {
			c.ILS.ShakeMin = cfg.ShakeMin
		}
		if cfg.ShakeMax > 0 {
			c.ILS.ShakeMax = cfg.ShakeMax
		}
		if cfg.MetropolisCriterion {
			c.ILS.MetropolisCriterion = cfg.MetropolisCriterion
		}
	}
}

func withVNS(cfg *vns.Configuration, timeLimitSeconds float64) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if timeLimitSeconds > 0 {
			c.VNS.TimeLimitSeconds = timeLimitSeconds
		} else if timeLimitSeconds == 0 {
			c.VNS.TimeLimitSeconds = math.MaxFloat64
		}

		if cfg.MaxIterations > 0 {
			c.VNS.MaxIterations = cfg.MaxIterations
		}
		if cfg.Rate > 0 {
			c.VNS.Rate = cfg.Rate
		}
	}
}
