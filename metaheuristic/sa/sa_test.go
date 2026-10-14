package sa_test

import (
	"context"
	"testing"

	"github.com/RKO-solver/rko-go/internal/testenv"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/sa"
	"github.com/RKO-solver/rko-go/metaheuristic/search"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
)

// TestSolveBaseline fails when a change alters the result for a fixed seed; if the change is intended, update want.
func TestSolveBaseline(t *testing.T) {
	want := []testenv.Result{
		{Cost: 5137, Hash: 0x55b0094fa85a5ed},
		{Cost: 5368, Hash: 0xa78ea0a9a0cbcb9c},
		{Cost: 5037, Hash: 0x170a34dd7582a919},
	}

	env := testenv.NewTSP(40, 1)
	newMetaheuristic := func(env problem.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) metaheuristic.Algorithm {
		// a low temperature and short inner loop make SA both accept and reject worse neighbours
		c := sa.DefaultConfig()
		c.Iterations = []int{10}
		c.TemperatureInitial = []float64{200}
		c.Alpha = []float64{0.05}
		return sa.New(env, rg, pool, l, sa.WithConfig(c), sa.WithSearch(search.NewDefault(env, pool, rg)))
	}
	ctx := context.Background()

	for i, w := range want {
		got := testenv.Run(t, ctx, env, uint64(i+1), 30, newMetaheuristic)
		assert.Equal(t, w, got, "seed %d: got %#v", i+1, got)
	}
}
