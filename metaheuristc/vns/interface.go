package vns

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (vns *VNS) SetRG(rg *random.Generator) {
	vns.RG = rg
	vns.search.SetRG(rg)
}

func (vns *VNS) Name() string {
	return name
}

func (vns *VNS) Solve(ctx context.Context) definition.Result {
	rko, elapsed := vns.solve(ctx, vns.solutionPool)

	return definition.Result{
		Solution:        vns.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

var _ definition.Solver = (*VNS)(nil)
