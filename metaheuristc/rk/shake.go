package rk

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/random"
)

type shakeType = int

const ShakeMax = swapNeighbour + 1

const (
	randomShake shakeType = iota
	invertShake
	swapRandom
	swapNeighbour
)

func Shake(solution *metaheuristc.RandomKeyValue, min float64, max float64, rg *random.Generator, env definition.Environment) {
	n := solution.RK.Len()
	var intensity int

	if min == max {
		intensity = int(float64(n)*min) + 1
	} else {
		intensity = int(float64(n)*rg.RangeFloat64(min, max)) + 1
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
