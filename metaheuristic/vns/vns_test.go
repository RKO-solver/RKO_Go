package vns_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/vns"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestSolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestSolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5156, Hash: 0x8a38ff1ff1028dc5},
		{Cost: 5119, Hash: 0xa5ed95fd88f14ac},
		{Cost: 5119, Hash: 0x48a0f47d831ca473},
	}

	env := testenv.NewTSP(40, 1)
	newSolver := func(env definition.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) definition.Solver {
		return vns.New(env, rg, pool, l)
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 30, newSolver)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
