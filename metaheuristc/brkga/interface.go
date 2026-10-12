package brkga

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
)

func (brkga *BRKGA) Name() string {
	return name
}

func (brkga *BRKGA) Solve(ctx context.Context) definition.Result {
	rko, elapsed := brkga.solve(ctx, brkga.solutionPool)

	return definition.Result{
		Solution: brkga.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*BRKGA)(nil)
