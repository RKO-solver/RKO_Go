package channel

import (
	"sync"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// solverInfo is one Register call sent through the update channel.
type solverInfo struct {
	localBest int
	local     int
	time      float64
	extra     string
}

// solverHistory holds one solver's unspilled records, its spill file and the last Register shown by the dashboard.
// All fields are guarded by information.mu.
type solverHistory struct {
	id      int // index in information.solvers
	name    string
	last    solverInfo    // the last Register, shown by the dashboard; kept apart because records is emptied on each spill
	hasLast bool          // false until the first Register of the run
	records []logger.Data // records not yet spilled to the file
	file    spill.File
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
	altScreen          bool // the dashboard is drawn in the terminal's alternate screen
	numLinesPool       int
	numVerboseMessages int
	pool               []poolInfo
	solvers            []*solverHistory
	extraMessages      [][]extraInfo
	workerMessages     []string
}

func (d *information) registerInfo(id int, info solverInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	s := d.solvers[id]
	s.last = info
	s.hasLast = true
	s.records = append(s.records, logger.Data{LocalCost: info.local, BestCost: info.localBest, Time: info.time})
	// Runs in the aggregator goroutine, so a full buffer is spilled here, under d.mu.
	// After a spill the slice is reused from [:0]; if spilling is off or fails, it keeps growing in memory.
	if len(s.records) == cap(s.records) && s.file.Write(s.records) {
		s.records = s.records[:0]
	}
}

func (d *information) registerVerbose(id int, extra extraInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.extraMessages[id] = append(d.extraMessages[id], extra)
}

// reset empties the recorded run and deletes its history files, applying c to every registered solver.
func (d *information) reset(c spill.Config) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for id, s := range d.solvers {
		s.reset(c)
		d.extraMessages[id] = make([]extraInfo, 0)
	}
	d.pool = make([]poolInfo, 0)
	d.workerMessages = make([]string, 0)
	d.altScreen = false
}

// reset deletes the spill file, applies c for this solver and empties the history; the caller holds information.mu.
func (s *solverHistory) reset(c spill.Config) {
	s.last = solverInfo{}
	s.hasLast = false
	c.Solver, c.ID = s.name, s.id
	s.file.Reset(c)
	s.records = make([]logger.Data, 0, spill.Capacity(c.Bytes)) // the capacity is the flush point
}
