package ga

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (ga *BRKGA) SetRG(rg *random.Generator) {
	ga.RG = rg
	ga.search.SetRG(ga.RG)
}

func (ga *BRKGA) Name() string {
	return nameBRKGA
}

func (ga *BRKGA) Solve(ctx context.Context) definition.Result {
	rko, elapsed := ga.solve(ctx, ga.solutionPool)

	return definition.Result{
		Solution: ga.env.Decode(rko.RK),
		Cost:     rko.Cost,
		Elapsed:  elapsed,
	}
}

var _ definition.Solver = (*BRKGA)(nil)
