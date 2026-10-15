package brkga_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/brkga"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestBRKGASolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestBRKGASolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5769, Hash: 0x2d01da0554d0772f},
		{Cost: 5669, Hash: 0x8801e49e1a75dd2a},
		{Cost: 5478, Hash: 0xb4d10fb578deff89},
	}

	env := testenv.NewTSP(40, 1)
	newMetaheuristic := func(env problem.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) metaheuristic.Algorithm {
		return brkga.New(env, rg, pool, l)
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 100, newMetaheuristic)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
