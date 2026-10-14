package config_test

import (
	"testing"
	"time"

	"github.com/RKO-solver/rko-go/config"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeParameterLists(t *testing.T) {
	src := "Parameters:\n  BRKGA:\n    EliteRatio: [0.1, 0.2]\n    PopulationSize: 50\n  SA:\n    QtdReheat:\n      - 1\n      - 2\n"
	for _, strict := range []bool{true, false} {
		cfg := mustDecode(t, src, strict)
		p := cfg.Parameters
		assert.Equal(t, metaheuristic.List[float64]{0.1, 0.2}, p.BRKGA.EliteRatio, "strict=%v", strict)
		assert.Equal(t, metaheuristic.List[int]{50}, p.BRKGA.PopulationSize, "strict=%v", strict)
		assert.Equal(t, metaheuristic.List[uint8]{1, 2}, p.SA.QtdReheat, "strict=%v", strict)
		assert.Equal(t, config.DefaultParameters().BRKGA.MutantRatio, p.BRKGA.MutantRatio, "strict=%v", strict)
	}
}

func TestDecodeParameterListsUnknownKey(t *testing.T) {
	src := "Parameters:\n  GA:\n    PopulationSize: [10, 20]\n    Bogus: 1\n"

	_, err := decodeOne(src, true)
	assert.ErrorIs(t, err, config.ErrUnknownField)

	cfg, err := decodeOne(src, false)
	require.NoError(t, err)
	assert.Equal(t, metaheuristic.List[int]{10, 20}, cfg.Parameters.GA.PopulationSize)
}

func TestDecodeParameterListWrongType(t *testing.T) {
	_, err := decodeOne("Parameters:\n  VNS:\n    Rate: [0.1, abc]\n", false)
	require.Error(t, err)
}

func TestNewSolverParameterListErrors(t *testing.T) {
	tests := []struct {
		name   string
		change func(p *config.Parameters)
		want   string
	}{
		{"duplicate", func(p *config.Parameters) { p.GA.PopulationSize = []int{10, 20, 10} }, "config: Parameters.GA.PopulationSize: invalid value: repeated value"},
		{"ILS shake", func(p *config.Parameters) { p.ILS.ShakeMin = []float64{0.1, 0.5}; p.ILS.ShakeMax = []float64{0.4, 0.6} }, "config: Parameters.ILS.ShakeMin: invalid value: 0.5 is greater than the maximum 0.4"},
		{"SA shake", func(p *config.Parameters) { p.SA.ShakeMin = []float64{0.9} }, "config: Parameters.SA.ShakeMin: invalid value: 0.9 is greater than the maximum 0.3"},
		{"LNS beta", func(p *config.Parameters) { p.LNS.BetaMax = []float64{0.001} }, "config: Parameters.LNS.BetaMin: invalid value: 0.01 is greater than the maximum 0.001"},
		{"SA goal duplicate", func(p *config.Parameters) { p.SA.TemperatureGoal = []float64{0.1, 0.1} }, "config: Parameters.SA.TemperatureGoal: invalid value: repeated value"},
		{"SA goal above initial", func(p *config.Parameters) { p.SA.TemperatureGoal = []float64{0.1, 2e7} }, "config: Parameters.SA.TemperatureGoal: invalid value: 2e+07 is not below the minimum TemperatureInitial 1e+07"},
		{"SA goal equal to initial", func(p *config.Parameters) { p.SA.TemperatureGoal = []float64{1e7} }, "config: Parameters.SA.TemperatureGoal: invalid value: 1e+07 is not below the minimum TemperatureInitial 1e+07"},
		{"empty list uses default", func(p *config.Parameters) { p.ILS.ShakeMax = nil; p.ILS.ShakeMin = []float64{0.04} }, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := silentConfig()
			cfg.Metaheuristics = []config.Metaheuristic{{Type: metaheuristic.GA}}
			cfg.Parameters = config.DefaultParameters()
			tt.change(cfg.Parameters)

			solver, log, err := cfg.NewSolver(testEnv{n: 30})
			if tt.want == "" {
				require.NoError(t, err)
				require.NoError(t, log.Close())
				return
			}
			assert.Nil(t, solver)
			assert.Nil(t, log)
			assert.ErrorIs(t, err, config.ErrInvalidValue)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestNewSolverTunedParameters(t *testing.T) {
	params := config.DefaultParameters()
	params.BRKGA.PopulationSize = []int{20, 30}
	params.GA.PopulationSize = []int{21, 10}
	params.SA.Alpha = []float64{0.1, 0.3}
	params.SA.Iterations = []int{3}
	params.ILS.ShakeMax = []float64{0.05, 0.1}
	params.VNS.Rate = []float64{0.1, 0.2}
	params.LNS.BetaMax = []float64{0.05, 0.1}
	before := cloneConfig(&config.Config{Parameters: params})

	cfg := silentConfig()
	cfg.MaxIterations = 20
	cfg.Metaheuristics = allMetaheuristics()
	cfg.Parameters = params

	solver, log := newSolver(t, cfg, config.WithSeed(7))
	solveWithin(t, solver, 30*time.Second)

	assert.Equal(t, before.Parameters, params)
	require.Len(t, log.Report(), len(cfg.Metaheuristics))
}
