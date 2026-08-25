package configuration

import (
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
)

func withMultiStart(cfg *multistart.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}
		if cfg.MaxIterations > 0 {
			c.MultiStart.MaxIterations = cfg.MaxIterations
		}
	}
}

func withBRKGA(cfg *ga.ConfigurationBRKGA) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
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

func withGA(cfg *ga.ConfigurationGA) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}

		if cfg.PopulationSize > 0 {
			c.GA.PopulationSize = cfg.PopulationSize
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

func withSA(cfg *sa.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}

		if cfg.MaxIterations > 0 {
			c.SA.MaxIterations = cfg.MaxIterations
		}
		if cfg.Alpha > 0 {
			c.SA.Alpha = cfg.Alpha
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

func withILS(cfg *ils.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
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

func withVNS(cfg *vns.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}

		if cfg.MaxIterations > 0 {
			c.VNS.MaxIterations = cfg.MaxIterations
		}
		if cfg.Rate > 0 {
			c.VNS.Rate = cfg.Rate
		}
	}
}

func withLNS(cfg *lns.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}

		if cfg.MaxIterations > 0 {
			c.LNS.MaxIterations = cfg.MaxIterations
		}
		if cfg.BetaMin > 0 {
			c.LNS.BetaMin = cfg.BetaMin
		}
		if cfg.BetaMax > 0 {
			c.LNS.BetaMax = cfg.BetaMax
		}
	}
}
