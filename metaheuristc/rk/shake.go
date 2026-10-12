package rk

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

type shakeType = int

// ShakeMax is the number of Shake moves.
const ShakeMax = swapNeighbour + 1

const (
	randomShake shakeType = iota
	invertShake
	swapRandom
	swapNeighbour
)

// Shake applies random moves to solution and recomputes its cost.
func Shake(solution *metaheuristc.RandomKeyValue, lo float64, hi float64, rg *random.Generator, env definition.Environment) {
	n := solution.RK.Len()
	var intensity int

	if lo == hi {
		intensity = int(float64(n)*lo) + 1
	} else {
		intensity = int(float64(n)*rg.RangeFloat64(lo, hi)) + 1
	}

	for k := 0; k < intensity; k++ {
		shakeId := rg.IntN(ShakeMax)

		i := rg.IntN(n)

		switch shakeId {
		case randomShake:
			solution.RK[i] = rg.Float64()
		case invertShake:
			solution.RK[i] = 1.0 - solution.RK[i]
		case swapRandom:
			j := rg.IntN(n)
			solution.RK[i], solution.RK[j] = solution.RK[j], solution.RK[i]
		case swapNeighbour:
			j := (i + 1) % n
			solution.RK[i], solution.RK[j] = solution.RK[j], solution.RK[i]
		}

	}

	solution.Cost = env.Cost(solution.RK)
}
