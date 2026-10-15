package multistart

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (m *MultiStart) Name() string {
	return name
}

func (m *MultiStart) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := m.solve(ctx, m.solutionPool)

	return metaheuristic.Result{
		Solution: m.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*MultiStart)(nil)
