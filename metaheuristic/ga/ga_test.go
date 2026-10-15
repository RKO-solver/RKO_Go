package ga_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/ga"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestGASolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestGASolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5657, Hash: 0x581cb2972964d9e0},
		{Cost: 5725, Hash: 0x35feddb5068a4c49},
		{Cost: 5027, Hash: 0x2d2355c54a483e},
	}

	env := testenv.NewTSP(40, 1)
	newMetaheuristic := func(env problem.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) metaheuristic.Algorithm {
		return ga.New(env, rg, pool, l)
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 100, newMetaheuristic)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
