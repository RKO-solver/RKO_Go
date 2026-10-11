package lns

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (lns *LNS) SetRG(rg *random.Generator) {
	lns.RG = rg
	lns.search.SetRG(rg)
}

func (lns *LNS) Name() string {
	return name
}

func (lns *LNS) Solve(ctx context.Context) definition.Result {
	rko, elapsed := lns.solve(ctx, lns.solutionPool)

	return definition.Result{
		Solution: lns.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*LNS)(nil)
