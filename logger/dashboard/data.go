package dashboard

import (
	"sync"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// metaheuristicHistory holds one metaheuristic's unspilled records, its spill file and the last Register shown by the dashboard.
// All fields are guarded by information.mu.
type metaheuristicHistory struct {
	id      int // index in information.metaheuristics
	name    string
	last    metaheuristicInfo // the last Register, shown by the dashboard; kept apart because records is emptied on each spill
	hasLast bool              // false until the first Register of the run
	records []logger.Sample   // records not yet spilled to the file
	events  []event           // the latest event of each key
	file    spill.File
	policy  *logger.Table // the tuner's learned policy; nil until Policy
}

type information struct {
	mu             sync.Mutex
	altScreen      bool // the dashboard is drawn in the terminal's alternate screen
	numLinesPool   int
	pool           []logger.PoolEntry
	metaheuristics []*metaheuristicHistory
	workers        []logger.WorkerResult
}

func (d *information) apply(m message) {
	switch m.t {
	case infoMessage:
		d.registerInfo(m.id, m.info)
	case eventMessage:
		d.registerEvent(m.id, m.event)
	}
}

func (d *information) registerInfo(id int, info metaheuristicInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	s := d.metaheuristics[id]
	s.last = info
	s.hasLast = true
	s.records = append(s.records, logger.Sample{Local: info.local, Best: info.best, Elapsed: info.elapsed})
	// Runs in the aggregator goroutine, so a full buffer is spilled here, under d.mu.
	// After a spill the slice is reused from [:0]; if spilling is off or fails, it keeps growing in memory.
	if len(s.records) == cap(s.records) && s.file.Write(s.records) {
		s.records = s.records[:0]
	}
}

func (d *information) setPolicy(id int, t logger.Table) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.metaheuristics[id].policy = &t
}

func (d *information) registerEvent(id int, e event) {
	d.mu.Lock()
	defer d.mu.Unlock()

	s := d.metaheuristics[id]
	for i := range s.events {
		if s.events[i].attr.Key == e.attr.Key {
			s.events[i] = e
			return
		}
	}
	s.events = append(s.events, e)
}

// reset empties the recorded run and deletes its history files, applying c to every registered metaheuristic.
func (d *information) reset(c spill.Config) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, s := range d.metaheuristics {
		s.reset(c)
	}
	d.pool = make([]logger.PoolEntry, 0, 50)
	d.workers = make([]logger.WorkerResult, 0)
	d.altScreen = false
}

// reset deletes the spill file, applies c for this metaheuristic and empties the history; the caller holds information.mu.
func (s *metaheuristicHistory) reset(c spill.Config) {
	s.last = metaheuristicInfo{}
	s.hasLast = false
	s.events = nil
	s.policy = nil
	c.Metaheuristic, c.ID = s.name, s.id
	s.file.Reset(c)
	s.records = make([]logger.Sample, 0, spill.Capacity(c.Bytes)) // the capacity is the flush point
}
