// Package solution manages pools of ranked candidate solutions.
package solution

import (
	"math"
	"slices"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/random"
)

// DefaultSize is the default maximum and initial pool size.
const DefaultSize = 10

//const similarlyMinimum = 0.8

// Pool is a thread-safe set of solutions sorted by ascending cost.
type Pool struct {
	mu        sync.Mutex
	solutions []*metaheuristc.RandomKeyValue
	best      atomic.Int64
	//similarity bool
	maxSize int
	logger  logger.Logger
}

// NewPoolSize creates a Pool of initialSize random solutions; maxSize <= 0 means no cap, else it is at least initialSize.
func NewPoolSize(env definition.Environment, rg *random.Generator, logger logger.Logger, maxSize int, initialSize int) *Pool {
	pool := &Pool{
		maxSize:   maxSize,
		logger:    logger,
		solutions: make([]*metaheuristc.RandomKeyValue, 0, max(maxSize, initialSize)),
	}

	if maxSize > 0 {
		pool.maxSize = max(initialSize, maxSize)
	}

	for range initialSize {
		key := rk.Generate(env, rg)
		cost := env.Cost(key)
		solution := &metaheuristc.RandomKeyValue{
			RK:   key,
			Cost: cost,
		}
		pool.solutions = append(pool.solutions, solution)
	}

	sort.Slice(pool.solutions, func(i, j int) bool { return pool.solutions[i].Cost < pool.solutions[j].Cost })
	pool.best.Store(math.MaxInt)
	if len(pool.solutions) > 0 {
		pool.best.Store(int64(pool.solutions[0].Cost))
	}
	return pool
}

// NewPool creates a Pool of DefaultSize random solutions capped at DefaultSize.
func NewPool(env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	return NewPoolSize(env, rg, logger, DefaultSize, DefaultSize)
}

// Add adds a copy of solution only if its cost is below the best; start sets the logged elapsed time.
func (p *Pool) Add(solution *metaheuristc.RandomKeyValue, start time.Time) {
	if int64(solution.Cost) >= p.best.Load() {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.solutions) == 0 {
		p.solutions = append(p.solutions, solution.Clone())
		p.best.Store(int64(solution.Cost))
		p.logger.AddSolutionPool(solution.Cost, time.Since(start))
		return
	}

	if p.solutions[0].Cost <= solution.Cost {
		return
	}

	//if p.solutions[0].Cost < solution.Cost {
	//	return
	//}
	//
	//if p.similarity && p.solutions[0].Cost == solution.Cost {
	//	if p.solutions[0].RK.Similarity(solution.RK) >= similarlyMinimum {
	//		return
	//	}
	//}

	p.solutions = slices.Insert(p.solutions, 0, solution.Clone())
	p.best.Store(int64(solution.Cost))
	p.logger.AddSolutionPool(solution.Cost, time.Since(start))
	if p.maxSize > 0 && len(p.solutions) >= p.maxSize {
		p.solutions = p.solutions[:p.maxSize]
	}

}

// Best returns a clone of the lowest-cost solution, or nil if empty.
func (p *Pool) Best() *metaheuristc.RandomKeyValue {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.solutions) == 0 {
		return nil
	}

	return p.solutions[0].Clone()
}

// Solution returns a clone of the solution at index (0 is the best).
func (p *Pool) Solution(index int) *metaheuristc.RandomKeyValue {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.solutions[index].Clone()
}

// BestCost returns the lowest cost, or math.MaxInt if empty.
func (p *Pool) BestCost() int {
	return int(p.best.Load())
}

// Len returns the number of solutions in the pool.
func (p *Pool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.solutions)
}
