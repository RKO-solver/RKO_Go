package channel

import (
	"sync"
)

type solverInfo struct {
	name      string
	id        int
	localBest int
	local     int
	time      float64
	extra     string
}

type poolInfo struct {
	cost int
	time float64
}

type extraInfo struct {
	message   string
	timeStamp float64
}
type information struct {
	mu                 sync.RWMutex
	init               bool
	numLinesPool       int
	numVerboseMessages int
	pool               []poolInfo
	solvers            [][]solverInfo
	extraMessages      [][]extraInfo
	workerMessages     []string
	previousLineCount  int
}

func (d *information) registerInfo(id int, info solverInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.solvers[id] = append(d.solvers[id], info)
}

func (d *information) registerVerbose(id int, extra extraInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.extraMessages[id] = append(d.extraMessages[id], extra)
}

// reset empties the recorded run, keeping one slot per registered solver.
func (d *information) reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for id := range d.solvers {
		d.solvers[id] = make([]solverInfo, 0)
		d.extraMessages[id] = make([]extraInfo, 0)
	}
	d.pool = make([]poolInfo, 0)
	d.workerMessages = make([]string, 0)
	d.init = false
	d.previousLineCount = 0
}
