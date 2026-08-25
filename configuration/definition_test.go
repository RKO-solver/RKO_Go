package configuration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RKO-solver/rko-go/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithGAYaml(t *testing.T) {
	yaml := `
GA:
  PopulationSize: 123
  CrossoverAlpha: 0.7
  MutationAlpha: 0.2
  MaxGenerations: 50
  MaxGenerationNoImprovement: 10
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ga.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlMHConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 123, cfg.GA.PopulationSize)
	assert.Equal(t, 0.7, cfg.GA.CrossoverAlpha)
	assert.Equal(t, 0.2, cfg.GA.MutationAlpha)
	assert.Equal(t, 50, cfg.GA.MaxGenerations)
	assert.Equal(t, 10, cfg.GA.MaxGenerationNoImprovement)
}

func TestWithBRKGAYaml(t *testing.T) {
	yaml := `
BRKGA:
  PopulationSize: 200
  EliteRatio: 0.2
  MutantRatio: 0.1
  CrossoverAlpha: 0.6
  MutationAlpha: 0.15
  MaxGenerations: 80
  MaxGenerationNoImprovement: 20
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "brkga.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlMHConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 200, cfg.BRKGA.PopulationSize)
	assert.Equal(t, 0.2, cfg.BRKGA.EliteRatio)
	assert.Equal(t, 0.1, cfg.BRKGA.MutantRatio)
	assert.Equal(t, 0.6, cfg.BRKGA.CrossoverAlpha)
	assert.Equal(t, 0.15, cfg.BRKGA.MutationAlpha)
	assert.Equal(t, 80, cfg.BRKGA.MaxGenerations)
	assert.Equal(t, 20, cfg.BRKGA.MaxGenerationNoImprovement)
}

func TestWithSAYaml(t *testing.T) {
	yaml := `
SA:
  MaxIterations: 500
  Alpha: 0.95
  TemperatureInitial: 1000
  TemperatureGoal: 10
  TemperatureReheat: 200
  ShakeMin: 2.0
  ShakeMax: 5.0
  QtdReheat: 3
  Iterations: 100
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "sa.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlMHConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 500, cfg.SA.MaxIterations)
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
VNS:
  MaxIterations: 77
  Rate: 0.33
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "vns.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlMHConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 77, cfg.VNS.MaxIterations)
	assert.Equal(t, 0.33, cfg.VNS.Rate)
}

func TestWithILSYaml(t *testing.T) {
	yaml := `
ILS:
  MaxIterations: 88
  ShakeMin: 2.0
  ShakeMax: 8.0
  MetropolisCriterion: true
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ils.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlMHConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 88, cfg.ILS.MaxIterations)
	assert.Equal(t, 2.0, cfg.ILS.ShakeMin)
	assert.Equal(t, 8.0, cfg.ILS.ShakeMax)
	assert.True(t, cfg.ILS.MetropolisCriterion)
}

func TestWithLNSYaml(t *testing.T) {
	yaml := `
LNS:
  MaxIterations: 88
  BetaMin: 2.0
  BetaMax: 8.0
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "lns.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlMHConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 88, cfg.LNS.MaxIterations)
	assert.Equal(t, 2.0, cfg.LNS.BetaMin)
	assert.Equal(t, 8.0, cfg.LNS.BetaMax)
}
