package ils_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/ils"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestSolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestSolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5658, Hash: 0x6a5e23a118450b30},
		{Cost: 5206, Hash: 0x248dccc3885af0d9},
		{Cost: 5477, Hash: 0xbedacc8a48716390},
	}

	env := testenv.NewTSP(40, 1)
	newMetaheuristic := func(env problem.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) metaheuristic.Algorithm {
		return ils.New(env, rg, pool, l)
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 100, newMetaheuristic)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
