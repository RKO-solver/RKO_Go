package brkga

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (brkga *BRKGA) Name() string {
	return name
}

func (brkga *BRKGA) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := brkga.solve(ctx, brkga.solutionPool)

	return metaheuristic.Result{
		Solution: brkga.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*BRKGA)(nil)
