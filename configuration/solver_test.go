package configuration_test

import (
	"math"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go"
	"github.com/RKO-solver/rko-go/configuration"
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testEnv is a cheap Environment with n keys.
type testEnv struct{ n int }

func (e testEnv) NumKeys() int { return e.n }
func (e testEnv) Cost(r definition.RandomKey) int {
	c := 0
	for i, v := range r {
		c += int(v*1000) * (i + 1)
	}
	return c
}
func (e testEnv) Decode(r definition.RandomKey) any { return r.Clone() }
func (e testEnv) SwapSearch() [][2]int              { return [][2]int{{0, e.n}} }

// solveWithin runs solver.Solve and fails the test if it does not return within d.
func solveWithin(t *testing.T, solver *rko.Solver, d time.Duration) {
	t.Helper()
	done := make(chan struct{})
	go func() { solver.Solve(); close(done) }()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatal("solver did not stop")
	}
}

func TestMaxIterationsStopsEveryMetaheuristic(t *testing.T) {
	const maxIterations = 5
	mhConfig := configuration.DefaultConfiguration()
	mhConfig.SA.Iterations = 10
	solvers := make([]configuration.Solver, 0)
	for _, mh := range []rko.MetaHeuristic{rko.MULTISTART, rko.SA, rko.GA, rko.VNS, rko.ILS, rko.BRKGA, rko.LNS} {
		solvers = append(solvers, configuration.Solver{MetaHeuristic: mh, Search: []search.Type{search.Swap}})
	}
	cfg := &configuration.SolverConfiguration{
		LoggerLevel:   logger.SILENT,
		LoggerType:    logger.PRINT,
		MaxIterations: maxIterations,
		Solvers:       solvers,
		Configuration: mhConfig,
	}
	solver, log := configuration.CreateSolverSeed("test", testEnv{n: 30}, cfg, 42)

	solveWithin(t, solver, 30*time.Second)

	report := log.GetReportData()
	require.Len(t, report, len(solvers))
	for _, info := range report {
		assert.Len(t, info.Performance, maxIterations, info.Name)
	}
}

func TestTimeLimitStopsBeforeMaxIterations(t *testing.T) {
	cfg := &configuration.SolverConfiguration{
		LoggerLevel:   logger.SILENT,
		LoggerType:    logger.PRINT,
		TimeLimit:     200 * time.Millisecond,
		MaxIterations: math.MaxInt,
		Solvers:       []configuration.Solver{{MetaHeuristic: rko.MULTISTART, Search: []search.Type{search.Swap}}},
	}
	solver, _ := configuration.CreateSolverSeed("test", testEnv{n: 30}, cfg, 42)
	solveWithin(t, solver, 10*time.Second)
}
