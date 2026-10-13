package brkga_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/brkga"
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
func run(t *testing.T, cfg *brkga.Config, deadline bool) string {
	t.Helper()

	env := tinyEnv{}
	l := printer.New(io.Discard, env.ProblemName(), nil)
	rg := random.New(1)
	pool := solution.NewPool(env, rg, l)
	solver := brkga.New(env, rg.Split(), pool, l, brkga.WithConfig(cfg))

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
		cfg := brkga.DefaultConfig()
		cfg.PopulationSize = metaheuristic.List[int]{10, 20}
		cfg.EliteRatio = metaheuristic.List[float64]{0.1, 0.3}
		cfg.MutantRatio = metaheuristic.List[float64]{0.05, 0.1}
		cfg.MaxGenerationNoImprovement = metaheuristic.List[int]{2, 5}
		assert.NotEmpty(t, run(t, cfg, deadline), "deadline=%v", deadline)
	}
}

func TestUntunedRunHasNoPolicy(t *testing.T) {
	assert.Empty(t, run(t, brkga.DefaultConfig(), false))
}

func TestNewPanicsOnInvalidEliteSize(t *testing.T) {
	for _, c := range []struct {
		size  metaheuristic.List[int]
		elite metaheuristic.List[float64]
	}{
		{metaheuristic.List[int]{10}, metaheuristic.List[float64]{0.05}},
		{metaheuristic.List[int]{10}, metaheuristic.List[float64]{1}},
		{metaheuristic.List[int]{100, 4}, metaheuristic.List[float64]{0.2}},
	} {
		cfg := brkga.DefaultConfig()
		cfg.PopulationSize, cfg.EliteRatio = c.size, c.elite
		env := tinyEnv{}
		l := printer.New(io.Discard, env.ProblemName(), nil)
		rg := random.New(1)
		pool := solution.NewPool(env, rg, l)
		assert.Panics(t, func() { brkga.New(env, rg, pool, l, brkga.WithConfig(cfg)) })
	}
}
