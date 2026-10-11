// Package definition defines the problem Environment, the Solver interface and random keys.
package definition

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/RKO-solver/rko-go/random"
)

// RandomKey is a solution encoded as keys in [0, 1).
type RandomKey []float64

// Solver is a metaheuristic that can be run by rko.Solver.
type Solver interface {
	// Solve runs the metaheuristic until it finishes or ctx is done and returns its Result.
	Solve(ctx context.Context) Result
	// Name returns the solver name.
	Name() string
	// SetRG sets the random generator of the solver.
	SetRG(rg *random.Generator)
	// Print prints the solver configuration.
	Print()
}

// Environment describes the problem to optimize; the user implements it.
type Environment interface {
	// NumKeys returns the number of keys in a solution.
	NumKeys() int
	// Cost returns the objective value of r; lower is better.
	Cost(r RandomKey) int
	// Decode converts r into the problem's own representation.
	Decode(r RandomKey) any
	// SwapSearch returns [start, end) key ranges; the Swap local search swaps keys within each range.
	SwapSearch() [][2]int
	// ProblemName returns the problem name used to label the logger and its output files.
	ProblemName() string
}

// Result is the outcome of one Solver.Solve call.
type Result struct {
	Solution any           // Solution decoded by Environment.Decode
	Cost     int           // Cost of the solution
	Elapsed  time.Duration // from the start of the run until Solve returned
}

// SortedIndex returns the indices that sort keys in ascending order (https://doi.org/10.48550/arXiv.2411.04293).
func (keys RandomKey) SortedIndex() []int {
	indices := make([]int, len(keys))
	for i := range indices {
		indices[i] = i
	}
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

// Equals reports whether keys and other hold exactly the same values.
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

// Similarity returns the cosine similarity of keys and key.
func (keys RandomKey) Similarity(key RandomKey) float64 {
	return dotMultiplier(keys, key) / (sizeVector(keys) * sizeVector(key))
}

func dotMultiplier(key1 RandomKey, key2 RandomKey) float64 {
	result := 0.0

	for i := range key1 {
		result += key1[i] * key2[i]
	}

	return result
}

func sizeVector(key RandomKey) float64 {
	result := 0.0

	for _, v := range key {
		result += v * v
	}

	return math.Sqrt(result)
}
