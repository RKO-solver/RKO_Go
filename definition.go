// Package rko runs several random-key metaheuristics in parallel over a shared solution pool.
package rko

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
)

// Solver runs several metaheuristics in parallel over a shared solution pool.
type Solver struct {
	l             logger.Logger
	env           definition.Environment
	solutionPool  *solution.Pool
	solvers       []definition.Solver
	timeLimit     time.Duration
	maxIterations int
}

// Solve runs the metaheuristics and returns the pool's best solution, decoded.
func (s *Solver) Solve(ctx context.Context) any {
	s.l.Reset()

	var wg sync.WaitGroup
	if s.maxIterations > 0 {
		ctx = definition.WithMaxIterations(ctx, s.maxIterations)
	}

	stop := func() {}
	if r, ok := s.l.(logger.Runner); ok {
		stop = r.Start()
	}

	ctxStartTime, cancel := definition.WithStartTime(ctx, s.timeLimit)
	defer cancel()
	for i, sv := range s.solvers {
		wg.Add(1)
		go metaheuristic.Worker(ctxStartTime, sv, i, s.l, &wg)
	}

	wg.Wait()
	stop()

	rk := s.solutionPool.Best()
	return s.env.Decode(rk.RK)
}

// Pool returns the pool shared by every running metaheuristic.
func (s *Solver) Pool() *solution.Pool {
	return s.solutionPool
}

// Print prints the configuration of every metaheuristic in the solver.
func (s *Solver) Print() {
	for i, sv := range s.solvers {
		fmt.Printf("(%d) ", i)
		sv.Print()
	}
}
