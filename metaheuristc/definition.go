// Package metaheuristc holds the types shared by the metaheuristics and the Worker that runs them.
package metaheuristc

import (
	"sort"

	"github.com/RKO-solver/rko-go/definition"
)

// RandomKeyValue is a random key with its cost.
type RandomKeyValue struct {
	RK   definition.RandomKey
	Cost int
}

// Clone returns a deep copy of rko.
func (rko *RandomKeyValue) Clone() *RandomKeyValue {
	sol := &RandomKeyValue{
		RK:   make(definition.RandomKey, len(rko.RK)),
		Cost: rko.Cost,
	}
	copy(sol.RK, rko.RK)

	return sol
}

// Sort sorts rkos by ascending cost.
func Sort(rkos []*RandomKeyValue) {
	sort.Slice(rkos, func(i, j int) bool { return rkos[i].Cost < rkos[j].Cost })
}
