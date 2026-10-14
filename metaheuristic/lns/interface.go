package lns

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (lns *LNS) Name() string {
	return name
}

func (lns *LNS) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := lns.solve(ctx, lns.solutionPool)

	return metaheuristic.Result{
		Solution: lns.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*LNS)(nil)
