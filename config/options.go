package config

import "github.com/RKO-solver/rko-go/random"

// Option customises NewSolver.
type Option func(*options)

type options struct {
	seed   uint64
	seeded bool
}

// WithSeed seeds the solver's random generator.
func WithSeed(seed uint64) Option { return func(o *options) { o.seed, o.seeded = seed, true } }

func (o *options) generator() *random.Generator {
	if o.seeded {
		return random.NewGeneratorSeed(o.seed)
	}
	return random.NewGenerator()
}
