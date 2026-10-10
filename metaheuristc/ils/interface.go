package ils

import (
	"context"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func (ils *ILS) SetRG(rg *random.Generator) {
	ils.RG = rg
	ils.search.SetRG(ils.RG)
}

func (ils *ILS) Name() string {
	return name
}

func (ils *ILS) Solve(ctx context.Context) definition.Result {
	rko, elapsed := ils.solve(ctx, ils.solutionPool)

	return definition.Result{
		Solution:        ils.env.Decode(rko.RK),
		Cost:            rko.Cost,
		TimeSpentSecond: elapsed,
	}
}

var _ definition.Solver = (*ILS)(nil)
