// Package definition provides core interfaces and types for defining optimization problems
// and solvers in the RKO. Users must implement the Environment interface
// for their specific problem to use the metaheuristics provided by the library.
package definition

import (
	"sort"

	"github.com/RKO-solver/rko-go/random"
)

// RandomKey represents a solution as a slice of float64 values (random keys).
type RandomKey []float64

// Solver is the interface that wraps the basic methods required for a metaheuristic solver.
type Solver interface {
	// Solve executes the metaheuristic and returns a Result.
	Solve() Result
	// Name returns the name of the solver.
	Name() string
	// SetRG sets the random number generator for the solver.
	SetRG(rg *random.Generator)
	// SetTimeLimitSecond sets the time limit for the solver in seconds.
	SetTimeLimitSecond(timeLimitSecond float64)
}

// Environment is the interface that must be implemented for a specific optimization problem.
// It defines how solutions are represented, evaluated, and decoded.
type Environment interface {
	// NumKeys returns the number of keys (variables) in the problem.
	NumKeys() int
	// Cost evaluates the cost (objective function) for a given solution.
	Cost(r RandomKey) int
	// Decode converts a solution from random keys to the problem's representation.
	Decode(r RandomKey) any
}

// Result holds the output of a solver, including the decoded solution, its cost, and the time spent.
type Result struct {
	Solution        any     // Decoded solution in the problem's representation
	Cost            int     // Cost (objective value) of the solution
	TimeSpentSecond float64 // Time spent to obtain the solution (in seconds)
}

// SortedIndex returns the indices that would sort the RandomKey in ascending order to be used by the user for decoding the problem (https://doi.org/10.48550/arXiv.2411.04293).
func (keys RandomKey) SortedIndex() []int {
	indices := make([]int, len(keys))
	for i := range indices {
		indices[i] = i
	}
	// Sort the indices based on the values in keys
	sort.Slice(indices, func(i, j int) bool {
		return keys[indices[i]] < keys[indices[j]]
	})
	return indices
}

// Len returns the number of elements in the RandomKey.
func (keys RandomKey) Len() int {
	return len(keys)
}

// Clone returns a copy of the RandomKey.
func (keys RandomKey) Clone() RandomKey {
	copyKeys := make(RandomKey, len(keys), cap(keys))
	copy(copyKeys, keys)
	return copyKeys
}

// Equals checks if two RandomKey slices are equal.
// Note: Be mindful of float precision if strict equality is not desired.
func (keys RandomKey) Equals(other RandomKey) bool {
	if keys == nil && other == nil {
		return true
	}
	if keys == nil || other == nil || len(keys) != len(other) {
		return false
	}
	for i := range keys {
		if keys[i] != other[i] {
			return false
		}
	}
	return true
}
