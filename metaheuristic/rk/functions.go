// Package rk generates and perturbs random keys.
package rk

import (
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/problem"
	"github.com/RKO-solver/rko-go/random"
)

// Generate returns a new random key for env.
func Generate(env problem.Environment, r *random.Generator) problem.RandomKey {
	numKeys := env.NumKeys()

	keys := make(problem.RandomKey, numKeys)

	for i := 0; i < numKeys; i++ {
		keys[i] = r.Float64()
	}

	return keys
}

// Reset fills randomKey with new random values in [0, 1).
func Reset(randomKey problem.RandomKey, rg *random.Generator) {
	for i := 0; i < len(randomKey); i++ {
		randomKey[i] = rg.Float64()
	}
}

// Copy copies src into dst; both must have the same length.
func Copy(dst, src *metaheuristic.RandomKeyValue) {
	copy(dst.RK, src.RK)
	dst.Cost = src.Cost
}
