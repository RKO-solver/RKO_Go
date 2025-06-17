package rk

import (
	"github.com/lucasmends/rko-go/definition"
	"github.com/lucasmends/rko-go/random"
)

/*
Generates a new randomKey
*/
func Generate(env definition.Environment, r *random.Generator) definition.RandomKey {
	numKeys := env.NumKeys()

	keys := make(definition.RandomKey, numKeys)

	for i := 0; i < numKeys; i++ {
		keys[i] = r.Float64()
	}

	return keys
}

/*
Generates random values for the randomKey given
*/
func Reset(randomKey definition.RandomKey, rg *random.Generator) {
	for i := 0; i < len(randomKey); i++ {
		randomKey[i] = rg.Float64()
	}
}
