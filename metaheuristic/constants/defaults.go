// Package constants holds the default parameters of the solver and the metaheuristics.
package constants

import "time"

// DefaultTimeLimit is the time limit of rko.New.
const DefaultTimeLimit = 150 * time.Second

// DefaultShakeMin is the default ILS ShakeMin and LNS BetaMin.
const DefaultShakeMin = 0.01

// DefaultShakeMax is the default ILS ShakeMax and LNS BetaMax.
const DefaultShakeMax = 0.05

// DefaultRate is the default VNS Rate.
const DefaultRate = 0.2

// DefaultAlphaSimulationAnnealing is the default SA Alpha.
const DefaultAlphaSimulationAnnealing = 0.0005

// DefaultTemperatureInitial is the default SA TemperatureInitial.
const DefaultTemperatureInitial = 10000000.0

// DefaultTemperatureGoal is the default SA TemperatureGoal.
const DefaultTemperatureGoal = 0.0000000000000001

// DefaultReheat is the default SA TemperatureReheat.
const DefaultReheat = 5.0

// DefaultShakeMinSimulationAnnealing is the default SA ShakeMin.
const DefaultShakeMinSimulationAnnealing = 0.1

// DefaultShakeMaxSimulationAnnealing is the default SA ShakeMax.
const DefaultShakeMaxSimulationAnnealing = 0.3

// DefaultPreheat is the default SA QtdReheat.
const DefaultPreheat = 0

// DefaultIterationsSimulationAnnealing is the default SA Iterations.
const DefaultIterationsSimulationAnnealing = 1000

// DefaultRestart is the default SA Restart.
const DefaultRestart = true

// DefaultPopulationSize is the default GA and BRKGA PopulationSize.
const DefaultPopulationSize = 200

// DefaultCrossoverAlpha is the default GA and BRKGA CrossoverAlpha.
const DefaultCrossoverAlpha = 0.95

// DefaultMutationAlpha is the default GA and BRKGA MutationAlpha.
const DefaultMutationAlpha = 0.005

// DefaultMaxGenerationNoImprovement is the default GA and BRKGA MaxGenerationNoImprovement.
const DefaultMaxGenerationNoImprovement = 100

// DefaultEliteRatio is the default BRKGA EliteRatio.
const DefaultEliteRatio = 0.2

// DefaultMutantRatio is the default BRKGA MutantRatio.
const DefaultMutantRatio = 0.05
