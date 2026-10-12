package sa

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
)

func (sa *SimulatedAnnealing) Name() string {
	return name
}

func (sa *SimulatedAnnealing) Solve(ctx context.Context) definition.Result {
	rko, elapsed := sa.solve(ctx, sa.solutionPool)

	return definition.Result{
		Solution: sa.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*SimulatedAnnealing)(nil)
