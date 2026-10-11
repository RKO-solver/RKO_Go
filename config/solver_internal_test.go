package config

import (
	"math"
	"testing"

	rko "github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedPoolSizes(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int
		initialSize int
		wantMax     int
		wantInitial int
	}{
		{"both zero", 0, 0, 10, 10},
		{"max above initial", 20, 5, 20, 5},
		{"max raised to initial", 3, 8, 8, 8},
		{"zero max follows larger initial", 0, 25, 25, 25},
		{"negative max is unlimited", -1, 0, -1, 10},
		{"large negative max is unlimited", -7, 30, -1, 30},
		{"MaxInt is unlimited", math.MaxInt, 4, -1, 4},
		{"negative initial uses default", 5, -3, 10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMax, gotInitial := seedPoolSizes(tt.maxSize, tt.initialSize)
			assert.Equal(t, tt.wantMax, gotMax)
			assert.Equal(t, tt.wantInitial, gotInitial)

			againMax, againInitial := seedPoolSizes(gotMax, gotInitial)
			assert.Equal(t, gotMax, againMax)
			assert.Equal(t, gotInitial, againInitial)
		})
	}
}

func TestNewSolverValidateOrder(t *testing.T) {
	cfg := &Config{
		LogLevel:       9,
		LogType:        7,
		Metaheuristics: []Metaheuristic{{Type: rko.GA}, {Type: rko.GRASP}, {Type: 99}},
		Parameters:     &Parameters{},
	}
	err := cfg.validate()
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownValue)
	assert.ErrorIs(t, err, ErrNotImplemented)
	assert.Equal(t, []string{
		"config: LogLevel: unknown value 9",
		"config: LogType: unknown value 7",
		"config: Metaheuristics[1]: Greedy Randomized Adaptive Search Procedure is not implemented",
		"config: Metaheuristics[2]: unknown value 99",
		"config: Parameters.BRKGA is nil",
		"config: Parameters.GA is nil",
		"config: Parameters.ILS is nil",
		"config: Parameters.SA is nil",
		"config: Parameters.VNS is nil",
		"config: Parameters.LNS is nil",
	}, splitLines(err.Error()))
}

func TestNewSolverValidateValid(t *testing.T) {
	cfg := &Config{LogType: logger.CHANNEL, Parameters: DefaultParameters()}
	assert.NoError(t, cfg.validate())
	cfg.LogType = logger.PRINT
	cfg.Parameters = nil
	assert.NoError(t, cfg.validate())
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}
