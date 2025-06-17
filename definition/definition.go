package definition

import (
	"github.com/lucasmends/rko-go/random"
	"sort"
)

type RandomKey []float64

type Solver interface {
	Solve() Result
	Name() string
	SetIdWorker(id int)
	SetRG(rg *random.Generator)
	SetTimeLimitSecond(timeLimitSecond float64)
}

type Environment interface {
	NumKeys() int
	Cost(r RandomKey) int
	Decode(r RandomKey) any
}

type Result struct {
	Solution        any
	Cost            int
	TimeSpentSecond float64
}

func (keys RandomKey) SortedIndex() []int {

	indices := make([]int, len(keys))

	for i := range indices {
		indices[i] = i
	}

	// Sort the indices based on the values in pieceKeys
	sort.Slice(indices, func(i, j int) bool {
		return keys[indices[i]] < keys[indices[j]]
	})

	return indices
}

func (keys RandomKey) Len() int {
	return len(keys)
}

func (keys RandomKey) Clone() RandomKey {
	copyKeys := make(RandomKey, len(keys), cap(keys))
	copy(copyKeys, keys)

	return copyKeys
}

func (keys RandomKey) Equals(other RandomKey) bool {
	if keys == nil && other == nil {
		return true
	}

	if keys == nil || other == nil || len(keys) != len(other) {
		return false
	}

	for i := range keys {
		if keys[i] != other[i] { // Again, be mindful of float precision if needed.
			return false
		}
	}
	return true
}
