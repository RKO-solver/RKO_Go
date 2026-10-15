package config

import (
	"math"

	"github.com/RKO-solver/rko-go/metaheuristic/solution"
)

func seedPoolSizes(maxSize, initialSize int) (int, int) {
	if initialSize < 1 {
		initialSize = solution.DefaultSize
	}
	if maxSize == 0 {
		maxSize = solution.DefaultSize
	}
	if maxSize < 0 || maxSize == math.MaxInt {
		return -1, initialSize
	}
	return max(maxSize, initialSize), initialSize
}
