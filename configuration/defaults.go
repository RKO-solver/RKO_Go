package configuration

import (
	"github.com/lucasmends/rko-go/metaheuristc/constants"
	"github.com/lucasmends/rko-go/metaheuristc/ga"
	"github.com/lucasmends/rko-go/metaheuristc/ils"
	"github.com/lucasmends/rko-go/metaheuristc/multistart"
	"github.com/lucasmends/rko-go/metaheuristc/sa"
	"github.com/lucasmends/rko-go/metaheuristc/vns"
)

func DefaultMultiStart() *multistart.Configuration {
	return &multistart.Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
	}
}

func DefaultBRKGA() *ga.ConfigurationBRKGA {
	return &ga.ConfigurationBRKGA{
		TimeLimitSeconds:           constants.DefaultTimeLimitSeconds,
		PopulationSize:             constants.DefaultPopulationSize,
		EliteRatio:                 constants.DefaultEliteRatio,
		MutantRatio:                constants.DefaultMutantRation,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerations:             constants.DefaultMaxGenerations,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}
}

func DefaultGA() *ga.ConfigurationGA {
	return &ga.ConfigurationGA{
		TimeLimitSeconds:           constants.DefaultTimeLimitSeconds,
		PopulationSize:             constants.DefaultPopulationSize,
		CrossoverAlpha:             constants.DefaultCrossoverAlpha,
		MutationAlpha:              constants.DefaultMutationAlpha,
		MaxGenerations:             constants.DefaultMaxGenerations,
		MaxGenerationNoImprovement: constants.DefaultMaxGenerationNoImprovement,
	}
}

func DefaultSA() *sa.Configuration {
	return &sa.Configuration{
		MaxIterations:      constants.DefaultMaxIterations,
		TimeLimitSeconds:   constants.DefaultTimeLimitSeconds,
		Alpha:              constants.DefaultAlphaSimulationAnnealing,
		ChangeImpact:       constants.DefaultImpact,
		TemperatureInitial: constants.DefaultTemperatureInitial,
		TemperatureGoal:    constants.DefaultTemperatureGoal,
		TemperatureReheat:  constants.DefaultReheat,
		ShakeMin:           constants.DefaultShakeMinSimulationAnnealing,
		ShakeMax:           constants.DefaultShakeMaxSimulationAnnealing,
		QtdReheat:          constants.DefaultPreheat,
		Iterations:         constants.DefaultIterationsSimulationAnnealing,
	}
}

func DefaultILS() *ils.Configuration {
	return &ils.Configuration{
		MaxIterations:       constants.DefaultMaxIterations,
		TimeLimitSeconds:    constants.DefaultTimeLimitSeconds,
		ShakeMin:            constants.DefaultShakeMin,
		ShakeMax:            constants.DefaultShakeMax,
		MetropolisCriterion: false,
	}
}

func DefaultVNS() *vns.Configuration {
	return &vns.Configuration{
		MaxIterations:    constants.DefaultMaxIterations,
		TimeLimitSeconds: constants.DefaultTimeLimitSeconds,
		Rate:             constants.DefaultRate,
	}
}
