package logger

import (
	"sync"
)

type solverInfo struct {
	name      string
	localBest int
	local     int
	timeStamp float64
	extra     string
}
type information struct {
	mu                sync.RWMutex
	init              bool
	solutionCost      []int
	solvers           [][]solverInfo
	extraMessages     [][]string
	previousLineCount int
}
