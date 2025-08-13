package configuration_test

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/lucasmends/rko-go/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithMultiStartYaml(t *testing.T) {
	yaml := `
TimeLimitSeconds: 10
MultiStart:
  MaxIterations: 42
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "multistart.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 42, cfg.MultiStart.MaxIterations)
	assert.Equal(t, 10.0, cfg.MultiStart.TimeLimitSeconds)
}

func TestWithMultiStartYaml_TimeLimit(t *testing.T) {
	yaml := `
TimeLimitSeconds: 15
MultiStart:
  MaxIterations: 99
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "multistart_time.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 99, cfg.MultiStart.MaxIterations)
	assert.Equal(t, 15.0, cfg.MultiStart.TimeLimitSeconds)
}

func TestWithMultiStartYaml_NoTimeLimit(t *testing.T) {
	yaml := `
MultiStart:
  MaxIterations: 77
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "multistart_notime.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 77, cfg.MultiStart.MaxIterations)
	assert.Equal(t, math.MaxFloat64, cfg.MultiStart.TimeLimitSeconds)
}

func TestWithGAYaml(t *testing.T) {
	yaml := `
GA:
  PopulationSize: 123
  ChildrenRatio: 0.3
  CrossoverAlpha: 0.7
  MutationAlpha: 0.2
  MaxGenerations: 50
  MaxGenerationNoImprovement: 10
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ga.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 123, cfg.GA.PopulationSize)
	assert.Equal(t, 0.3, cfg.GA.ChildrenRatio)
	assert.Equal(t, 0.7, cfg.GA.CrossoverAlpha)
	assert.Equal(t, 0.2, cfg.GA.MutationAlpha)
	assert.Equal(t, 50, cfg.GA.MaxGenerations)
	assert.Equal(t, 10, cfg.GA.MaxGenerationNoImprovement)
}

func TestWithGAYaml_TimeLimit(t *testing.T) {
	yaml := `
TimeLimitSeconds: 22
GA:
  PopulationSize: 10
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ga_time.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 10, cfg.GA.PopulationSize)
	assert.Equal(t, 22.0, cfg.GA.TimeLimitSeconds)
}

func TestWithGAYaml_NoTimeLimit(t *testing.T) {
	yaml := `
GA:
  PopulationSize: 11
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ga_notime.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 11, cfg.GA.PopulationSize)
	assert.Equal(t, math.MaxFloat64, cfg.GA.TimeLimitSeconds)
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
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 200, cfg.BRKGA.PopulationSize)
	assert.Equal(t, 0.2, cfg.BRKGA.EliteRatio)
	assert.Equal(t, 0.1, cfg.BRKGA.MutantRatio)
	assert.Equal(t, 0.6, cfg.BRKGA.CrossoverAlpha)
	assert.Equal(t, 0.15, cfg.BRKGA.MutationAlpha)
	assert.Equal(t, 80, cfg.BRKGA.MaxGenerations)
	assert.Equal(t, 20, cfg.BRKGA.MaxGenerationNoImprovement)
}

func TestWithBRKGAYaml_TimeLimit(t *testing.T) {
	yaml := `
TimeLimitSeconds: 33
BRKGA:
  PopulationSize: 20
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "brkga_time.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 20, cfg.BRKGA.PopulationSize)
	assert.Equal(t, 33.0, cfg.BRKGA.TimeLimitSeconds)
}

func TestWithBRKGAYaml_NoTimeLimit(t *testing.T) {
	yaml := `
BRKGA:
  PopulationSize: 21
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "brkga_notime.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 21, cfg.BRKGA.PopulationSize)
	assert.Equal(t, math.MaxFloat64, cfg.BRKGA.TimeLimitSeconds)
}

func TestWithSAYaml(t *testing.T) {
	yaml := `
SA:
  MaxIterations: 500
  Alpha: 0.95
  ChangeImpact: 0.5
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
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 500, cfg.SA.MaxIterations)
	assert.Equal(t, 0.95, cfg.SA.Alpha)
	assert.Equal(t, 0.5, cfg.SA.ChangeImpact)
	assert.Equal(t, 1000.0, cfg.SA.TemperatureInitial)
	assert.Equal(t, 10.0, cfg.SA.TemperatureGoal)
	assert.Equal(t, 200.0, cfg.SA.TemperatureReheat)
	assert.Equal(t, 2.0, cfg.SA.ShakeMin)
	assert.Equal(t, 5.0, cfg.SA.ShakeMax)
	assert.Equal(t, uint8(3), cfg.SA.QtdReheat)
	assert.Equal(t, 100, cfg.SA.Iterations)
}

func TestWithSAYaml_TimeLimit(t *testing.T) {
	yaml := `
TimeLimitSeconds: 44
SA:
  MaxIterations: 30
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "sa_time.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 30, cfg.SA.MaxIterations)
	assert.Equal(t, 44.0, cfg.SA.TimeLimitSeconds)
}

func TestWithSAYaml_NoTimeLimit(t *testing.T) {
	yaml := `
SA:
  MaxIterations: 31
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "sa_notime.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 31, cfg.SA.MaxIterations)
	assert.Equal(t, math.MaxFloat64, cfg.SA.TimeLimitSeconds)
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
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 77, cfg.VNS.MaxIterations)
	assert.Equal(t, 0.33, cfg.VNS.Rate)
}

func TestWithVNSYaml_TimeLimit(t *testing.T) {
	yaml := `
TimeLimitSeconds: 55
VNS:
  MaxIterations: 40
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "vns_time.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 40, cfg.VNS.MaxIterations)
	assert.Equal(t, 55.0, cfg.VNS.TimeLimitSeconds)
}

func TestWithVNSYaml_NoTimeLimit(t *testing.T) {
	yaml := `
VNS:
  MaxIterations: 41
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "vns_notime.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 41, cfg.VNS.MaxIterations)
	assert.Equal(t, math.MaxFloat64, cfg.VNS.TimeLimitSeconds)
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
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 88, cfg.ILS.MaxIterations)
	assert.Equal(t, 2.0, cfg.ILS.ShakeMin)
	assert.Equal(t, 8.0, cfg.ILS.ShakeMax)
	assert.True(t, cfg.ILS.MetropolisCriterion)
}

func TestWithILSYaml_TimeLimit(t *testing.T) {
	yaml := `
TimeLimitSeconds: 66
ILS:
  MaxIterations: 50
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ils_time.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 50, cfg.ILS.MaxIterations)
	assert.Equal(t, 66.0, cfg.ILS.TimeLimitSeconds)
}

func TestWithILSYaml_NoTimeLimit(t *testing.T) {
	yaml := `
ILS:
  MaxIterations: 51
`
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "ils_notime.yaml")
	require.NoError(t, os.WriteFile(file, []byte(yaml), 0644))
	cfg, err := configuration.CreateYamlConfiguration(file)
	require.NoError(t, err)
	assert.Equal(t, 51, cfg.ILS.MaxIterations)
	assert.Equal(t, math.MaxFloat64, cfg.ILS.TimeLimitSeconds)
}
