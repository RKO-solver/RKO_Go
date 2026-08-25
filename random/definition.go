// Package random provides pseudo-random number generation helpers.
package random

import (
	"math/rand/v2"
	"sync"
	"time"
)

// Generator wraps a PCG-based pseudo-random source.
type Generator struct {
	rand *rand.Rand
}

var (
	instance *Generator
	once     sync.Once
)

// NewGenerator returns a Generator seeded from the current time.
func NewGenerator() *Generator {
	currentTimeSeed := uint64(time.Now().UnixNano())
	source := rand.NewPCG(currentTimeSeed, currentTimeSeed+1)

	return &Generator{
		rand: rand.New(source),
	}
}

// GetGlobalInstance returns the singleton Generator, created on first call.
func GetGlobalInstance() *Generator {
	once.Do(func() {
		instance = NewGenerator()
	})

	return instance
}

// NewGeneratorSeed returns a Generator seeded with seed, for reproducible runs.
func NewGeneratorSeed(seed uint64) *Generator {
	source := rand.NewPCG(seed, seed+1)

	return &Generator{
		rand: rand.New(source),
	}
}

// Float64 returns a pseudo-random number in [0.0, 1.0).
func (g *Generator) Float64() float64 {
	return g.rand.Float64()
}

// Float32 returns a pseudo-random number in [0.0, 1.0).
func (g *Generator) Float32() float32 {
	return g.rand.Float32()
}

// IntN returns a pseudo-random int in [0, n).
func (g *Generator) IntN(n int) int {
	return g.rand.IntN(n)
}

// RangeInts returns numElem pseudo-random ints in [0, maxInt).
func (g *Generator) RangeInts(maxInt, numElem int) []int {
	values := make([]int, numElem)
	for i := 0; i < numElem; i++ {
		values[i] = g.rand.IntN(maxInt)
	}
	return values
}

// Permutation returns a pseudo-random permutation of [0, n).
func (g *Generator) Permutation(n int) []int {
	return g.rand.Perm(n)
}

// RangeFloat64 returns a pseudo-random float64 in [min, max).
func (g *Generator) RangeFloat64(min, max float64) float64 {
	return min + g.rand.Float64()*(max-min)
}

// RangeInt returns a pseudo-random int in [min, max).
func (g *Generator) RangeInt(min, max int) int {
	return min + g.rand.IntN(max-min)
}
