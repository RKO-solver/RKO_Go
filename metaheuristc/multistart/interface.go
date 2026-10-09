package multistart

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

// SetRG sets the random generator of the solver and of its local search.
func (m *MultiStart) SetRG(rg *random.Generator) {
	m.RG = rg
	m.search.SetRG(m.RG)
}

// Name returns the solver name.
func (m *MultiStart) Name() string {
	return name
}

// Solve runs MultiStart until ctx is done or the iteration limit is reached.
func (m *MultiStart) Solve(ctx context.Context) definition.Result {
	rko, elapsed := m.solve(ctx, m.solutionPool)

	return definition.Result{
		Solution:        m.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

// Compile-time check that *MultiStart implements definition.Solver.
var _ definition.Solver = (*MultiStart)(nil)
