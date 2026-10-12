package ga

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
)

func (ga *GA) Name() string {
	return nameGA
}

func (ga *GA) Solve(ctx context.Context) definition.Result {
	rko, elapsed := ga.solve(ctx, ga.solutionPool)

	return definition.Result{
		Solution: ga.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*GA)(nil)
