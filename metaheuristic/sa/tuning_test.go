package sa_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/sa"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tinyEnv is a cheap Environment with 100 keys.
type tinyEnv struct{}

func (tinyEnv) NumKeys() int        { return 100 }
func (tinyEnv) ProblemName() string { return "tiny" }
func (tinyEnv) Cost(r definition.RandomKey) int {
	c := 0
	for i, v := range r {
		c += int(v*1000) * (i + 1)
	}
	return c
}
func (tinyEnv) Decode(r definition.RandomKey) any { return r.Clone() }
func (tinyEnv) SwapSearch() [][2]int              { return [][2]int{{0, 100}} }

// run solves tinyEnv with cfg for 40 iterations and returns the printed policies.
func run(t *testing.T, cfg *sa.Config, deadline bool) string {
	t.Helper()

	env := tinyEnv{}
	l := printer.New(io.Discard, env.ProblemName(), nil)
	rg := random.New(1)
	pool := solution.NewPool(env, rg, l)
	solver := sa.New(env, rg.Split(), pool, l, sa.WithConfig(cfg))

	ctx := definition.WithMaxIterations(context.Background(), 40)
	if deadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}
	ctx = definition.WithStartTime(ctx, time.Now())
	res := solver.Solve(ctx)

	assert.Equal(t, env.Cost(res.Solution.(definition.RandomKey)), res.Cost)
	assert.LessOrEqual(t, pool.BestCost(), res.Cost)

	var buf bytes.Buffer
	require.NoError(t, l.PrintPolicies(&buf))
	return buf.String()
}

func TestTunedRun(t *testing.T) {
	for _, deadline := range []bool{true, false} {
		cfg := sa.DefaultConfig()
		cfg.Alpha = metaheuristic.List[float64]{0.05, 0.1}
		cfg.TemperatureInitial = metaheuristic.List[float64]{100, 200}
		cfg.TemperatureReheat = metaheuristic.List[float64]{5, 50}
		cfg.QtdReheat = metaheuristic.List[uint8]{0, 2}
		cfg.Iterations = metaheuristic.List[int]{2, 4}
		assert.NotEmpty(t, run(t, cfg, deadline), "deadline=%v", deadline)
	}
}

func TestUntunedRunHasNoPolicy(t *testing.T) {
	cfg := sa.DefaultConfig()
	cfg.Iterations = metaheuristic.List[int]{3}
	assert.Empty(t, run(t, cfg, false))
}
