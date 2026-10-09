package search

import (
	"context"
	"math"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const defaultRho = 0.5
const defaultMu = 0.02

func blend(keyA definition.RandomKey, keyB definition.RandomKey, factor bool, rho float64, mu float64, rg *random.Generator) definition.RandomKey {
	blendedKey := make(definition.RandomKey, len(keyA))

	for i := range len(keyA) {
		if rg.Float64() < mu {
			blendedKey[i] = rg.Float64()
			continue
		}
		if rg.Float64() < rho {
			blendedKey[i] = keyA[i]
			continue
		}
		if factor {
			blendedKey[i] = keyB[i]
			continue
		}
		blendedKey[i] = 1.0 - keyB[i]
	}

	return blendedKey
}

func computeBlend(keyA *metaheuristc.RandomKeyValue, keyB *metaheuristc.RandomKeyValue, factor bool, env definition.Environment, rg *random.Generator) *metaheuristc.RandomKeyValue {
	auxRK := blend(keyA.RK, keyB.RK, factor, defaultRho, defaultMu, rg)

	return &metaheuristc.RandomKeyValue{
		RK:   auxRK,
		Cost: env.Cost(auxRK),
	}
}

func rank(x1, x2, p *metaheuristc.RandomKeyValue) (*metaheuristc.RandomKeyValue, *metaheuristc.RandomKeyValue, *metaheuristc.RandomKeyValue) {
	switch {
	case p.Cost < x1.Cost:
		return p, x1, x2
	case p.Cost < x2.Cost:
		return x1, p, x2
	default:
		return x1, x2, p
	}
}

type nelderMeadLocalSearch struct {
	environment   definition.Environment
	solutionPool  *solution.Pool
	rg            *random.Generator
	maxIterations int
}

func (n nelderMeadLocalSearch) String() string {
	return GetSearchString(Nelder)
}

func (n *nelderMeadLocalSearch) SetRG(rg *random.Generator) {
	n.rg = rg
}

func (n nelderMeadLocalSearch) Search(ctx context.Context, rko *metaheuristc.RandomKeyValue) {
	nelderMeadSearch(ctx, rko, n.maxIterations, n.environment, n.solutionPool, n.rg)
}

func CreateNelderMeadLocalSearch(environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local {
	maxIterations := int(float64(environment.NumKeys()) * math.Exp(-2))
	return &nelderMeadLocalSearch{environment, solutionPool, rg, maxIterations}
}

func nelderMeadSearch(ctx context.Context, rko *metaheuristc.RandomKeyValue, maxIterations int, env definition.Environment, solutionPool *solution.Pool, rg *random.Generator) {
	var x1, x2, x3, x0, xR, xAux *metaheuristc.RandomKeyValue

	poolSize := solutionPool.Size()
	x1 = solutionPool.GetSolution(rg.IntN(poolSize))
	x2 = solutionPool.GetSolution(rg.IntN(poolSize))

	if x2.Cost < x1.Cost {
		x1, x2 = x2, x1
	}
	x1, x2, x3 = rank(x1, x2, rko.Clone())

	for i := 0; i < maxIterations; i++ {
		if ctx.Err() != nil {
			break
		}

		x0 = computeBlend(x1, x2, true, env, rg)
		xR = computeBlend(x0, x3, false, env, rg)

		if xR.Cost < x1.Cost {
			xAux = computeBlend(xR, x0, false, env, rg)
			if xAux.Cost < xR.Cost {
				x1, x2, x3 = xAux, x1, x2
			} else {
				x1, x2, x3 = xR, x1, x2
			}
			continue
		}

		if xR.Cost < x2.Cost {
			x2, x3 = xR, x2
			continue
		}

		if xR.Cost < x3.Cost {
			xAux = computeBlend(xR, x0, true, env, rg)
			if xAux.Cost < xR.Cost {
				x1, x2, x3 = rank(x1, x2, xAux)
				continue
			}
		} else {
			xAux = computeBlend(x0, x3, true, env, rg)
			if xAux.Cost < x3.Cost {
				x1, x2, x3 = rank(x1, x2, xAux)
				continue
			}
		}

		// shrink
		x2 = computeBlend(x1, x2, true, env, rg)
		x3 = computeBlend(x1, x3, true, env, rg)
		if x2.Cost < x1.Cost {
			x1, x2 = x2, x1
		}
		x1, x2, x3 = rank(x1, x2, x3)

	}

	copy(rko.RK, x1.RK)
	rko.Cost = x1.Cost
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*nelderMeadLocalSearch)(nil)
