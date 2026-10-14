package rko_test

import (
	"context"
	"io"
	"testing"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/stretchr/testify/require"
)

func newLog(t *testing.T, env *testenv.TSP) *printer.Log {
	t.Helper()
	log := printer.New(io.Discard, env.ProblemName(), &printer.Options{Level: logger.LevelSilent})
	t.Cleanup(func() { require.NoError(t, log.Close()) })
	return log
}

func TestNewErrors(t *testing.T) {
	env := testenv.NewTSP(10, 1)
	log := newLog(t, env)

	tests := []struct {
		name string
		opts []rko.Option
		want error
	}{
		{"none", nil, rko.ErrNoMetaheuristics},
		{"not implemented", []rko.Option{rko.WithMetaheuristics(metaheuristic.GA, metaheuristic.GRASP)}, rko.ErrNotImplemented},
		{"unknown", []rko.Option{rko.WithMetaheuristics(metaheuristic.Type(99))}, rko.ErrUnknownMetaheuristic},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solver, err := rko.New(env, log, tt.opts...)
			require.ErrorIs(t, err, tt.want)
			require.Nil(t, solver)
		})
	}
}

func TestNewFromMetaheuristicsEmpty(t *testing.T) {
	env := testenv.NewTSP(10, 1)
	solver, err := rko.NewFromMetaheuristics(env, newLog(t, env), nil, nil, rko.WithMaxIterations(1))
	require.ErrorIs(t, err, rko.ErrNoMetaheuristics)
	require.Nil(t, solver)
}

func TestWithMetaheuristicsAppends(t *testing.T) {
	env := testenv.NewTSP(10, 1)
	log := newLog(t, env)
	solver, err := rko.New(env, log, rko.WithMetaheuristics(metaheuristic.GA), rko.WithMetaheuristics(metaheuristic.SA), rko.WithMaxIterations(1))
	require.NoError(t, err)

	solver.Solve(context.Background())
	require.Len(t, log.Report(), 2)
}
