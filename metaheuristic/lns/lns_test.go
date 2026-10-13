package lns_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/lns"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestSolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestSolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5463, Hash: 0x572384ee227e26d8},
		{Cost: 5685, Hash: 0x3d2cf298d971c867},
		{Cost: 5631, Hash: 0xf1e2e2c132eb1d1c},
	}

	env := testenv.NewTSP(40, 1)
	newSolver := func(env definition.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) definition.Solver {
		return lns.New(env, rg, pool, l)
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 100, newSolver)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
