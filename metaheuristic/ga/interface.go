package ga

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (ga *GA) Name() string {
	return name
}

func (ga *GA) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := ga.solve(ctx, ga.solutionPool)

	return metaheuristic.Result{
		Solution: ga.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*GA)(nil)
