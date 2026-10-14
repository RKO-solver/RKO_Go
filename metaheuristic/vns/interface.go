package vns

import (
	"context"

	"github.com/RKO-solver/rko-go/metaheuristic"
)

func (vns *VNS) Name() string {
	return name
}

func (vns *VNS) Solve(ctx context.Context) metaheuristic.Result {
	rko, elapsed := vns.solve(ctx, vns.solutionPool)

	return metaheuristic.Result{
		Solution: vns.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ metaheuristic.Algorithm = (*VNS)(nil)
