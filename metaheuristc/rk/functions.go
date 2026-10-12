// Package rk generates and perturbs random keys.
package rk

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

// Generate returns a new random key for env.
func Generate(env definition.Environment, r *random.Generator) definition.RandomKey {
	numKeys := env.NumKeys()

	keys := make(definition.RandomKey, numKeys)

	for i := 0; i < numKeys; i++ {
		keys[i] = r.Float64()
	}

	return keys
}

// Reset fills randomKey with new random values in [0, 1).
func Reset(randomKey definition.RandomKey, rg *random.Generator) {
	for i := 0; i < len(randomKey); i++ {
		randomKey[i] = rg.Float64()
	}
}

func Copy(dst, src *metaheuristc.RandomKeyValue) {
	copy(dst.RK, src.RK)
	dst.Cost = src.Cost
}
