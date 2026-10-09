package configuration

import (
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
)

// withBRKGA copies the positive fields of cfg into the BRKGA defaults.
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
		if cfg.MaxGenerationNoImprovement > 0 {
			c.BRKGA.MaxGenerationNoImprovement = cfg.MaxGenerationNoImprovement
		}
	}
}

// withGA copies the positive fields of cfg into the GA defaults.
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
		if cfg.MaxGenerationNoImprovement > 0 {
			c.GA.MaxGenerationNoImprovement = cfg.MaxGenerationNoImprovement
		}

	}
}

// withSA copies the positive fields of cfg into the SA defaults.
func withSA(cfg *sa.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
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

// withILS copies the positive fields of cfg into the ILS defaults; MetropolisCriterion can only be turned on.
func withILS(cfg *ils.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
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

// withVNS copies the positive fields of cfg into the VNS defaults.
func withVNS(cfg *vns.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}

		if cfg.Rate > 0 {
			c.VNS.Rate = cfg.Rate
		}
	}
}

// withLNS copies the positive fields of cfg into the LNS defaults.
func withLNS(cfg *lns.Configuration) Option {
	return func(c *MetaheuristicsConfiguration) {
		if cfg == nil {
			return
		}

		if cfg.BetaMin > 0 {
			c.LNS.BetaMin = cfg.BetaMin
		}
		if cfg.BetaMax > 0 {
			c.LNS.BetaMax = cfg.BetaMax
		}
	}
}
