package constants

import "math"

const DefaultTimeLimitSeconds = float64(150)
const DefaultMaxIterations = math.MaxInt
const DefaultShakeMin = 0.01
const DefaultShakeMax = 0.05
const DefaultRate = 0.2

const DefaultAlphaSimulationAnnealing = 0.0005
const DefaultImpact = 0.20
const DefaultTemperatureInitial = 10000000.0
const DefaultTemperatureGoal = 0.0000000000000001
const DefaultReheat = 5.0
const DefaultShakeMinSimulationAnnealing = 0.1
const DefaultShakeMaxSimulationAnnealing = 0.3
const DefaultPreheat = 0
const DefaultIterationsSimulationAnnealing = 1000

const DefaultPopulationSize = 200
const DefaultCrossoverAlpha = 0.95
const DefaultMutationAlpha = 0.005
const DefaultMaxGenerations = math.MaxInt
const DefaultMaxGenerationNoImprovement = 100

const DefaultEliteRatio = 0.2
const DefaultMutantRation = 0.05
