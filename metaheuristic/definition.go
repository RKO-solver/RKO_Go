// Package metaheuristic holds the types shared by the metaheuristics and the Worker that runs them.
package metaheuristic

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/RKO-solver/rko-go/problem"
)

// RandomKeyValue is a random key with its cost.
type RandomKeyValue struct {
	RK   problem.RandomKey
	Cost int
}

// Clone returns a deep copy of rko.
func (rko *RandomKeyValue) Clone() *RandomKeyValue {
	sol := &RandomKeyValue{
		RK:   make(problem.RandomKey, len(rko.RK)),
		Cost: rko.Cost,
	}
	copy(sol.RK, rko.RK)

	return sol
}

// NewRandomKeyValue returns numKeys zero keys with an unevaluated (MaxInt) cost.
func NewRandomKeyValue(numKeys int) *RandomKeyValue {
	return &RandomKeyValue{
		RK:   make(problem.RandomKey, numKeys),
		Cost: math.MaxInt,
	}
}

// Sort sorts rkos by ascending cost.
func Sort(rkos []*RandomKeyValue) {
	sort.Slice(rkos, func(i, j int) bool { return rkos[i].Cost < rkos[j].Cost })
}

// LocalSearch is a local search that improves a solution in place.
type LocalSearch interface {
	// Search improves rko in place, stopping early when ctx is done.
	Search(ctx context.Context, rko *RandomKeyValue)
	// String returns the search name.
	String() string
}

// Print prints a metaheuristic name and its local search.
func Print(name string, s LocalSearch) {
	fmt.Printf("%s:\n", name)
	fmt.Printf("    Local Search:\n       %s\n", s.String())
}
