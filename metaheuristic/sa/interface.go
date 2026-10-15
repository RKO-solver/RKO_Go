package sa

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (sa *SimulatedAnnealing) Name() string {
	return name
}

func (sa *SimulatedAnnealing) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := sa.solve(ctx, sa.solutionPool)

	return metaheuristic.Result{
		Solution: sa.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*SimulatedAnnealing)(nil)
