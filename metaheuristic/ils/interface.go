package ils

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (ils *ILS) Name() string {
	return name
}

func (ils *ILS) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := ils.solve(ctx, ils.solutionPool)

	return metaheuristic.Result{
		Solution: ils.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*ILS)(nil)
