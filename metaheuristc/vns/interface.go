package vns

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
)

func (vns *VNS) Name() string {
	return name
}

func (vns *VNS) Solve(ctx context.Context) definition.Result {
	rko, elapsed := vns.solve(ctx, vns.solutionPool)

	return definition.Result{
		Solution: vns.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*VNS)(nil)
