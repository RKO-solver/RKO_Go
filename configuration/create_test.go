package configuration_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/configuration"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeConfig writes content to a temporary config.yaml and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(file, []byte(content), 0644))
	return file
}

// loadConfig parses content with CreateSolverConfiguration, failing the test on error.
func loadConfig(t *testing.T, content string) *configuration.SolverConfiguration {
	t.Helper()
	cfg, err := configuration.CreateSolverConfiguration(writeConfig(t, content))
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotNil(t, cfg.Configuration)
	return cfg
}

func TestCreateSolverConfigurationFull(t *testing.T) {
	yaml := `
LogLevel: VERBOSE
LogType: CHANNEL
TimeLimitSeconds: 60
Metaheuristics:
  - "ILS=Swap, Mirror"
  - "BRKGA"
  - "SA=Farey"
Configuration:
  MultiStart:
    MaxIterations: 11
  BRKGA:
    PopulationSize: 300
    EliteRatio: 0.25
    MutantRatio: 0.1
    CrossoverAlpha: 0.6
    MutationAlpha: 0.15
    MaxGenerations: 80
    MaxGenerationNoImprovement: 20
  GA:
    PopulationSize: 123
    CrossoverAlpha: 0.7
    MutationAlpha: 0.2
    MaxGenerations: 50
    MaxGenerationNoImprovement: 10
  ILS:
    MaxIterations: 88
    ShakeMin: 2.0
    ShakeMax: 8.0
    MetropolisCriterion: true
  SA:
    Alpha: 0.95
    TemperatureInitial: 1000
    TemperatureGoal: 10
    TemperatureReheat: 200
    ShakeMin: 2.0
    ShakeMax: 5.0
    QtdReheat: 3
    Iterations: 100
  VNS:
    MaxIterations: 77
    Rate: 0.33
  LNS:
    MaxIterations: 99
    BetaMin: 0.5
    BetaMax: 0.9
`
	cfg := loadConfig(t, yaml)

	assert.Equal(t, logger.VERBOSE, cfg.LoggerLevel)
	assert.Equal(t, logger.CHANNEL, cfg.LoggerType)
	assert.Equal(t, 60*time.Second, cfg.TimeLimit)
	assert.Equal(t, []configuration.Solver{
		{MetaHeuristic: rko.ILS, Search: []search.Type{search.Swap, search.Mirror}},
		{MetaHeuristic: rko.BRKGA, Search: []search.Type{}},
		{MetaHeuristic: rko.SA, Search: []search.Type{search.Farey}},
	}, cfg.Solvers)

	mh := cfg.Configuration
	assert.Equal(t, 11, mh.MultiStart.MaxIterations)

	assert.Equal(t, 300, mh.BRKGA.PopulationSize)
	assert.Equal(t, 0.25, mh.BRKGA.EliteRatio)
	assert.Equal(t, 0.1, mh.BRKGA.MutantRatio)
	assert.Equal(t, 0.6, mh.BRKGA.CrossoverAlpha)
	assert.Equal(t, 0.15, mh.BRKGA.MutationAlpha)
	assert.Equal(t, 80, mh.BRKGA.MaxGenerations)
	assert.Equal(t, 20, mh.BRKGA.MaxGenerationNoImprovement)

	assert.Equal(t, 123, mh.GA.PopulationSize)
	assert.Equal(t, 0.7, mh.GA.CrossoverAlpha)
	assert.Equal(t, 0.2, mh.GA.MutationAlpha)
	assert.Equal(t, 50, mh.GA.MaxGenerations)
	assert.Equal(t, 10, mh.GA.MaxGenerationNoImprovement)

	assert.Equal(t, 88, mh.ILS.MaxIterations)
	assert.Equal(t, 2.0, mh.ILS.ShakeMin)
	assert.Equal(t, 8.0, mh.ILS.ShakeMax)
	assert.True(t, mh.ILS.MetropolisCriterion)

	assert.Equal(t, 0.95, mh.SA.Alpha)
	assert.Equal(t, 1000.0, mh.SA.TemperatureInitial)
	assert.Equal(t, 10.0, mh.SA.TemperatureGoal)
	assert.Equal(t, 200.0, mh.SA.TemperatureReheat)
	assert.Equal(t, 2.0, mh.SA.ShakeMin)
	assert.Equal(t, 5.0, mh.SA.ShakeMax)
	assert.Equal(t, uint8(3), mh.SA.QtdReheat)
	assert.Equal(t, 100, mh.SA.Iterations)

	assert.Equal(t, 77, mh.VNS.MaxIterations)
	assert.Equal(t, 0.33, mh.VNS.Rate)

	assert.Equal(t, 99, mh.LNS.MaxIterations)
	assert.Equal(t, 0.5, mh.LNS.BetaMin)
	assert.Equal(t, 0.9, mh.LNS.BetaMax)
}

func TestWithMultiStartYaml(t *testing.T) {
	yaml := `
Configuration:
  MultiStart:
    MaxIterations: 42
`
	cfg := loadConfig(t, yaml)
	assert.Equal(t, 42, cfg.Configuration.MultiStart.MaxIterations)
}

func TestWithGAYaml(t *testing.T) {
	yaml := `
Configuration:
  GA:
    PopulationSize: 123
    CrossoverAlpha: 0.7
    MutationAlpha: 0.2
    MaxGenerations: 50
    MaxGenerationNoImprovement: 10
`
	cfg := loadConfig(t, yaml).Configuration
	assert.Equal(t, 123, cfg.GA.PopulationSize)
	assert.Equal(t, 0.7, cfg.GA.CrossoverAlpha)
	assert.Equal(t, 0.2, cfg.GA.MutationAlpha)
	assert.Equal(t, 50, cfg.GA.MaxGenerations)
	assert.Equal(t, 10, cfg.GA.MaxGenerationNoImprovement)
}

func TestWithBRKGAYaml(t *testing.T) {
	yaml := `
Configuration:
  BRKGA:
    PopulationSize: 300
    EliteRatio: 0.25
    MutantRatio: 0.1
    CrossoverAlpha: 0.6
    MutationAlpha: 0.15
    MaxGenerations: 80
    MaxGenerationNoImprovement: 20
`
	cfg := loadConfig(t, yaml).Configuration
	assert.Equal(t, 300, cfg.BRKGA.PopulationSize)
	assert.Equal(t, 0.25, cfg.BRKGA.EliteRatio)
	assert.Equal(t, 0.1, cfg.BRKGA.MutantRatio)
	assert.Equal(t, 0.6, cfg.BRKGA.CrossoverAlpha)
	assert.Equal(t, 0.15, cfg.BRKGA.MutationAlpha)
	assert.Equal(t, 80, cfg.BRKGA.MaxGenerations)
	assert.Equal(t, 20, cfg.BRKGA.MaxGenerationNoImprovement)
}

func TestWithSAYaml(t *testing.T) {
	yaml := `
Configuration:
  SA:
    Alpha: 0.95
    TemperatureInitial: 1000
    TemperatureGoal: 10
    TemperatureReheat: 200
    ShakeMin: 2.0
    ShakeMax: 5.0
    QtdReheat: 3
    Iterations: 100
`
	cfg := loadConfig(t, yaml).Configuration
	assert.Equal(t, 0.95, cfg.SA.Alpha)
	assert.Equal(t, 1000.0, cfg.SA.TemperatureInitial)
	assert.Equal(t, 10.0, cfg.SA.TemperatureGoal)
	assert.Equal(t, 200.0, cfg.SA.TemperatureReheat)
	assert.Equal(t, 2.0, cfg.SA.ShakeMin)
	assert.Equal(t, 5.0, cfg.SA.ShakeMax)
	assert.Equal(t, uint8(3), cfg.SA.QtdReheat)
	assert.Equal(t, 100, cfg.SA.Iterations)
}

func TestWithVNSYaml(t *testing.T) {
	yaml := `
Configuration:
  VNS:
    MaxIterations: 77
    Rate: 0.33
`
	cfg := loadConfig(t, yaml).Configuration
	assert.Equal(t, 77, cfg.VNS.MaxIterations)
	assert.Equal(t, 0.33, cfg.VNS.Rate)
}

func TestWithILSYaml(t *testing.T) {
	yaml := `
Configuration:
  ILS:
    MaxIterations: 88
    ShakeMin: 2.0
    ShakeMax: 8.0
    MetropolisCriterion: true
`
	cfg := loadConfig(t, yaml).Configuration
	assert.Equal(t, 88, cfg.ILS.MaxIterations)
	assert.Equal(t, 2.0, cfg.ILS.ShakeMin)
	assert.Equal(t, 8.0, cfg.ILS.ShakeMax)
	assert.True(t, cfg.ILS.MetropolisCriterion)
}

func TestWithLNSYaml(t *testing.T) {
	yaml := `
Configuration:
  LNS:
    MaxIterations: 88
    BetaMin: 2.0
    BetaMax: 8.0
`
	cfg := loadConfig(t, yaml).Configuration
	assert.Equal(t, 88, cfg.LNS.MaxIterations)
	assert.Equal(t, 2.0, cfg.LNS.BetaMin)
	assert.Equal(t, 8.0, cfg.LNS.BetaMax)
}

func TestPartialSectionKeepsDefaults(t *testing.T) {
	yaml := `
Configuration:
  SA:
    Alpha: 0.9
    Iterations: 7
  ILS:
    ShakeMax: 0.4
`
	cfg := loadConfig(t, yaml).Configuration
	defaults := configuration.DefaultConfiguration()

	expectedSA := *defaults.SA
	expectedSA.Alpha = 0.9
	expectedSA.Iterations = 7
	assert.Equal(t, expectedSA, *cfg.SA)

	expectedILS := *defaults.ILS
	expectedILS.ShakeMax = 0.4
	assert.Equal(t, expectedILS, *cfg.ILS)

	assert.Equal(t, defaults.MultiStart, cfg.MultiStart)
	assert.Equal(t, defaults.BRKGA, cfg.BRKGA)
	assert.Equal(t, defaults.GA, cfg.GA)
	assert.Equal(t, defaults.VNS, cfg.VNS)
	assert.Equal(t, defaults.LNS, cfg.LNS)
}

func TestNonPositiveValuesKeepDefaults(t *testing.T) {
	yaml := `
Configuration:
  VNS:
    MaxIterations: 0
    Rate: -1
  ILS:
    MetropolisCriterion: false
`
	cfg := loadConfig(t, yaml).Configuration
	defaults := configuration.DefaultConfiguration()

	assert.Equal(t, defaults.VNS, cfg.VNS)
	assert.Equal(t, defaults.ILS, cfg.ILS)
}

func TestMissingSectionsUseDefaults(t *testing.T) {
	yaml := `
LogLevel: SILENT
TimeLimitSeconds: 5
Metaheuristics:
  - "GA"
`
	cfg := loadConfig(t, yaml)
	assert.Equal(t, configuration.DefaultConfiguration(), cfg.Configuration)
}

func TestTopLevelSectionsAreIgnored(t *testing.T) {
	yaml := `
GA:
  PopulationSize: 123
VNS:
  Rate: 0.33
`
	cfg := loadConfig(t, yaml)
	assert.Equal(t, configuration.DefaultConfiguration(), cfg.Configuration)
}

func TestEmptyFileUsesDefaults(t *testing.T) {
	cfg := loadConfig(t, "")

	assert.Equal(t, logger.DefaultLogLevel, cfg.LoggerLevel)
	assert.Equal(t, logger.PRINT, cfg.LoggerType)
	assert.Equal(t, time.Duration(0), cfg.TimeLimit)
	assert.Empty(t, cfg.Solvers)
	assert.Equal(t, configuration.DefaultConfiguration(), cfg.Configuration)
}

func TestCreateSolverConfigurationMissingFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "does-not-exist.yaml")
	cfg, err := configuration.CreateSolverConfiguration(file)
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestCreateSolverConfigurationInvalidYaml(t *testing.T) {
	file := writeConfig(t, "LogLevel: [INFO\nConfiguration: :::")
	cfg, err := configuration.CreateSolverConfiguration(file)
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to unmarshal yaml")
}

func TestTimeLimitSecondsConvertedToDuration(t *testing.T) {
	cfg := loadConfig(t, "TimeLimitSeconds: 90")
	assert.Equal(t, 90*time.Second, cfg.TimeLimit)
}

func TestLogLevelParsing(t *testing.T) {
	tests := []struct {
		value    string
		expected logger.Level
	}{
		{"SILENT", logger.SILENT},
		{"silent", logger.SILENT},
		{"INFO", logger.INFO},
		{"Info", logger.INFO},
		{"VERBOSE", logger.VERBOSE},
		{"verbose", logger.VERBOSE},
		{"DEBUG", logger.INFO},
		{`""`, logger.INFO},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			cfg := loadConfig(t, "LogLevel: "+tt.value)
			assert.Equal(t, tt.expected, cfg.LoggerLevel)
		})
	}
}

func TestLogTypeParsing(t *testing.T) {
	tests := []struct {
		value    string
		expected logger.LogType
	}{
		{"CHANNEL", logger.CHANNEL},
		{"channel", logger.CHANNEL},
		{"PRINT", logger.PRINT},
		{"Print", logger.PRINT},
		{"FILE", logger.PRINT},
		{`""`, logger.PRINT},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			cfg := loadConfig(t, "LogType: "+tt.value)
			assert.Equal(t, tt.expected, cfg.LoggerType)
		})
	}
}

func TestKeysAreCaseSensitive(t *testing.T) {
	yaml := `
loglevel: SILENT
logtype: CHANNEL
timelimitseconds: 10
`
	cfg := loadConfig(t, yaml)
	assert.Equal(t, logger.DefaultLogLevel, cfg.LoggerLevel)
	assert.Equal(t, logger.PRINT, cfg.LoggerType)
	assert.Equal(t, time.Duration(0), cfg.TimeLimit)
}

func TestMetaheuristicsParsing(t *testing.T) {
	tests := []struct {
		name     string
		entries  []string
		expected []configuration.Solver
	}{
		{
			name:    "name without equals gives empty search",
			entries: []string{"VNS"},
			expected: []configuration.Solver{
				{MetaHeuristic: rko.VNS, Search: []search.Type{}},
			},
		},
		{
			name:    "whitespace is trimmed",
			entries: []string{"  ils  =  swap ,  MIRROR  , Nelder "},
			expected: []configuration.Solver{
				{MetaHeuristic: rko.ILS, Search: []search.Type{search.Swap, search.Mirror, search.Nelder}},
			},
		},
		{
			name:    "names are case-insensitive",
			entries: []string{"multistart", "Lns=farey"},
			expected: []configuration.Solver{
				{MetaHeuristic: rko.MULTISTART, Search: []search.Type{}},
				{MetaHeuristic: rko.LNS, Search: []search.Type{search.Farey}},
			},
		},
		{
			name:    "unknown names are dropped",
			entries: []string{"TABU=Swap", "GA", "", "NOPE"},
			expected: []configuration.Solver{
				{MetaHeuristic: rko.GA, Search: []search.Type{}},
			},
		},
		{
			name:    "unknown searches and RVND are filtered out",
			entries: []string{"SA=Swap, Unknown, RVND, rvnd, Mirror"},
			expected: []configuration.Solver{
				{MetaHeuristic: rko.SA, Search: []search.Type{search.Swap, search.Mirror}},
			},
		},
		{
			name:    "only invalid searches leave an empty search",
			entries: []string{"BRKGA=RVND, Bogus", "VNS="},
			expected: []configuration.Solver{
				{MetaHeuristic: rko.BRKGA, Search: []search.Type{}},
				{MetaHeuristic: rko.VNS, Search: []search.Type{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yaml := "Metaheuristics:\n"
			for _, e := range tt.entries {
				yaml += "  - \"" + e + "\"\n"
			}
			cfg := loadConfig(t, yaml)
			assert.Equal(t, tt.expected, cfg.Solvers)
		})
	}
}
