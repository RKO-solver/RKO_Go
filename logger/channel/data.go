package channel

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

func (d *information) registerInfo(id int, info solverInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.solvers[id] = append(d.solvers[id], info)
}

func (d *information) registerVerbose(id int, message string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.extraMessages[id] = append(d.extraMessages[id], message)
}
