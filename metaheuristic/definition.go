// Package metaheuristic holds the types shared by the metaheuristics and the Worker that runs them.
package metaheuristic

import (
	"math"
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

// NewRandomKeyValue returns numKeys zero keys with an unevaluated (MaxInt) cost.
func NewRandomKeyValue(numKeys int) *RandomKeyValue {
	return &RandomKeyValue{
		RK:   make(definition.RandomKey, numKeys),
		Cost: math.MaxInt,
	}
}

// Sort sorts rkos by ascending cost.
func Sort(rkos []*RandomKeyValue) {
	sort.Slice(rkos, func(i, j int) bool { return rkos[i].Cost < rkos[j].Cost })
}
