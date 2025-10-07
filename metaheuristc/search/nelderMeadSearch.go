package search

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

const defaultRho = 0.5
const defaultMu = 0.02
const maxIterations = 5

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

type nelderMeadLocalSearch struct {
	environment  definition.Environment
	solutionPool *solution.Pool
	rg           *random.Generator
}

func (n nelderMeadLocalSearch) SetRG(rg *random.Generator) {
	n.rg = rg
}

func (n nelderMeadLocalSearch) Search(rko *metaheuristc.RandomKeyValue) {
	nelderMeadSearch(rko, n.environment, n.solutionPool, n.rg)
	if rko.Cost < n.solutionPool.BestSolutionCost() {
		n.solutionPool.AddSolution(rko.Clone())
	}
}

func CreateNelderMeadLocalSearch(environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local {
	return nelderMeadLocalSearch{environment, solutionPool, rg}
}

func nelderMeadSearch(rko *metaheuristc.RandomKeyValue, env definition.Environment, solutionPool *solution.Pool, rg *random.Generator) {
	var x1, x2, x3, x0, xR, xAux *metaheuristc.RandomKeyValue

	poolSize := solutionPool.Size()
	firstElement := rg.IntN(poolSize)
	x1 = solutionPool.GetSolution(firstElement)
	secondElement := rg.IntN(poolSize)
	for secondElement == firstElement {
		secondElement = rg.IntN(poolSize)
	}
	x2 = solutionPool.GetSolution(secondElement)

	if x2.Cost < x1.Cost {
		xAux = x1
		x1 = x2
		x2 = xAux
	}

	x3 = rko.Clone()
	if x3.Cost < x1.Cost {
		xAux = x1
		x1 = x3
		x3 = x2
		x2 = xAux
	} else if x3.Cost < x2.Cost {
		xAux = x2
		x2 = x3
		x3 = xAux
	}

	for i := 0; i < maxIterations; i++ {
		x0 = computeBlend(x1, x2, true, env, rg)
		xR = computeBlend(x0, x3, false, env, rg)

		if xR.Cost < x1.Cost {
			xAux = computeBlend(xR, x0, false, env, rg)

			x3 = x2
			x2 = x1
			if xAux.Cost < xR.Cost {
				x1 = xAux
			} else {
				x1 = xR
			}
			continue
		}

		if xR.Cost < x2.Cost {
			x3 = x2
			x2 = xR
			continue
		}

		if xR.Cost < x3.Cost {
			xAux = computeBlend(xR, x0, true, env, rg)
			if xAux.Cost < xR.Cost {
				x3 = xAux
				continue
			}
		} else {
			xAux = computeBlend(x0, x3, true, env, rg)
			if xAux.Cost < x3.Cost {
				x3 = xAux
				continue
			}
		}

		// shrink
		x2 = computeBlend(x1, x2, true, env, rg)
		x3 = computeBlend(x1, x3, true, env, rg)

		if x2.Cost < x1.Cost {
			xAux = x1
			x1 = x2
			x2 = xAux
		}
		if x3.Cost < x1.Cost {
			xAux = x1
			x3 = x2
			x1 = x3
			x2 = xAux
		} else if x3.Cost < x2.Cost {
			xAux = x2
			x2 = x3
			x3 = xAux
		}
	}

	copy(rko.RK, x1.RK)
	rko.Cost = x1.Cost
}

// --- The Compile-Time Check ---
// This line "tells" the compiler to verify that *MyProcessor implements DataProcessor.
// If it doesn't, the code will not compile.
var _ Local = (*nelderMeadLocalSearch)(nil)
