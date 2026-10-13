// Package testenv provides a deterministic TSP Environment for the metaheuristic tests.
package testenv

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"io"
	"math"
	"testing"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/printer"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/random"
)

// TSP is a symmetric Euclidean TSP; a key vector decodes to the tour of its sorted indices.
type TSP struct{ dist [][]int }

// NewTSP returns a TSP of n cities placed at random from seed.
func NewTSP(n int, seed uint64) *TSP {
	rg := random.New(seed)
	x, y := make([]float64, n), make([]float64, n)
	for i := range n {
		x[i], y[i] = rg.Float64()*1000, rg.Float64()*1000
	}

	dist := make([][]int, n)
	for i := range n {
		dist[i] = make([]int, n)
		for j := range n {
			dist[i][j] = int(math.Hypot(x[i]-x[j], y[i]-y[j]))
		}
	}

	return &TSP{dist: dist}
}

func (t *TSP) NumKeys() int { return len(t.dist) }

func (t *TSP) Cost(r definition.RandomKey) int {
	tour := r.SortedIndex()
	cost := 0
	for i := range tour {
		cost += t.dist[tour[i]][tour[(i+1)%len(tour)]]
	}
	return cost
}

func (t *TSP) Decode(r definition.RandomKey) any { return r.Clone() }

func (t *TSP) SwapSearch() [][2]int { return [][2]int{{0, len(t.dist)}} }

func (t *TSP) ProblemName() string { return "tsp" }

var _ definition.Environment = (*TSP)(nil)

// Result is the outcome of one seeded run: the cost and a hash of the returned keys.
type Result struct {
	Cost int
	Hash uint64
}

// Factory builds the solver under test.
type Factory func(env definition.Environment, rg *random.Generator, pool *solution.Pool, l logger.Logger) definition.Solver

// Run solves env with seed for iterations and checks the result against the cost function and the pool.
func Run(t *testing.T, ctx context.Context, env *TSP, seed uint64, iterations int, f Factory) Result {
	t.Helper()

	l := printer.New(io.Discard, env.ProblemName(), nil)
	rg := random.New(seed)
	pool := solution.NewPool(env, rg, l)
	solver := f(env, rg.Split(), pool, l)

	ctx = definition.WithMaxIterations(ctx, iterations)
	ctx, cancel := definition.WithStartTime(ctx, 0)
	defer cancel()
	res := solver.Solve(ctx)

	keys := res.Solution.(definition.RandomKey)
	if c := env.Cost(keys); c != res.Cost {
		t.Errorf("seed %d: returned cost %d, but its keys cost %d", seed, res.Cost, c)
	}
	if best := pool.BestCost(); best > res.Cost {
		t.Errorf("seed %d: pool best %d is worse than the returned cost %d", seed, best, res.Cost)
	}

	return Result{Cost: res.Cost, Hash: Hash(keys)}
}

// Hash returns an FNV-1a hash of the bits of r.
func Hash(r definition.RandomKey) uint64 {
	h := fnv.New64a()
	buf := make([]byte, 0, 8*len(r))
	for _, v := range r {
		buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(v))
	}
	h.Write(buf)
	return h.Sum64()
}
