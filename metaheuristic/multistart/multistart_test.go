package multistart_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic/multistart"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestSolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestSolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5509, Hash: 0x7c52a1da9d48ac21},
		{Cost: 6180, Hash: 0xf73915b78196fd2f},
		{Cost: 5587, Hash: 0x94087f3eff2870c0},
	}

	env := testenv.NewTSP(40, 1)
	newSolver := func(env definition.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) definition.Solver {
		return multistart.New(env, rg, pool, l)
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 50, newSolver)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
