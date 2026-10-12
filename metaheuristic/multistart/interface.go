package multistart

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
)

func (m *MultiStart) Name() string {
	return name
}

func (m *MultiStart) Solve(ctx context.Context) definition.Result {
	rko, elapsed := m.solve(ctx, m.solutionPool)

	return definition.Result{
		Solution: m.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*MultiStart)(nil)
