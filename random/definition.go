// Package random provides pseudo-random number generation helpers.
package random

import "math/rand/v2"

// Generator is a PCG-based random source, not safe for concurrent use; call Split per goroutine.
type Generator struct {
	rand *rand.Rand
}

// NewRandom returns a randomly seeded Generator.
func NewRandom() *Generator {
	return New(rand.Uint64())
}

// New returns a Generator that yields the same sequence for the same seed.
func New(seed uint64) *Generator {
	source := rand.NewPCG(seed, seed+1)

	return &Generator{
		rand: rand.New(source),
	}
}

// Split returns a new Generator seeded with two values drawn from g.
func (g *Generator) Split() *Generator {
	return &Generator{rand: rand.New(rand.NewPCG(g.rand.Uint64(), g.rand.Uint64()))}
}

// Float64 returns a pseudo-random number in [0.0, 1.0).
func (g *Generator) Float64() float64 {
	return g.rand.Float64()
}

// IntN returns a pseudo-random int in [0, n).
func (g *Generator) IntN(n int) int {
	return g.rand.IntN(n)
}

// Perm returns a pseudo-random permutation of [0, n).
func (g *Generator) Perm(n int) []int {
	return g.rand.Perm(n)
}

// RangeFloat64 returns a pseudo-random float64 in [lo, hi).
func (g *Generator) RangeFloat64(lo, hi float64) float64 {
	return lo + g.rand.Float64()*(hi-lo)
}

// RangeInt returns a pseudo-random int in [lo, hi); panics if hi <= lo.
func (g *Generator) RangeInt(lo, hi int) int {
	return lo + g.rand.IntN(hi-lo)
}
