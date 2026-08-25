// Package solution manages pools of ranked candidate solutions.
package solution

import (
	"math"
	"sort"
	"sync"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/random"
)

const defaultMaxSize = 10

// Pool is a thread-safe, size-bounded collection of ranked solutions.
type Pool struct {
	mu        sync.RWMutex
	solutions []*metaheuristc.RandomKeyValue
	limited   bool
	maxSize   int
	logger    logger.Logger
}

var (
	instance *Pool
	once     sync.Once
)

// GetGlobalInstance returns the singleton Pool, created on the first call.
func GetGlobalInstance(env definition.Environment, logger logger.Logger, rg *random.Generator) *Pool {
	once.Do(func() {
		instance = NewDefaultPool(env, rg, logger)
	})

	return instance
}

// NewPool creates a Pool seeded with initialSize random solutions.
func NewPool(maxSize int, initialSize int, env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	pool := &Pool{
		maxSize:   maxSize,
		logger:    logger,
		limited:   maxSize > 1 && maxSize < math.MaxInt,
		solutions: make([]*metaheuristc.RandomKeyValue, 0, initialSize),
	}

	if pool.limited && initialSize > maxSize {
		pool.maxSize = initialSize
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
	return pool
}

// NewDefaultPool creates a Pool capped at 10 random solutions.
func NewDefaultPool(env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	return NewPool(defaultMaxSize, defaultMaxSize, env, rg, logger)
}

// NewDefaultPoolUnlimited creates an unbounded Pool with 10 random solutions.
func NewDefaultPoolUnlimited(env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	return NewPool(-1, defaultMaxSize, env, rg, logger)
}

// AddSolution inserts solution into the pool if it is not the worst.
func (p *Pool) AddSolution(solution *metaheuristc.RandomKeyValue, time float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.solutions) == 0 {
		p.solutions = append(p.solutions, solution)
		p.logger.AddSolutionPool(solution.Cost, time)
		return
	}

	// np better than the worst ignore
	if p.solutions[len(p.solutions)-1].Cost < solution.Cost {
		return
	}

	p.logger.AddSolutionPool(solution.Cost, time)
	p.solutions = append(p.solutions, solution)
	sort.Slice(p.solutions, func(i, j int) bool { return p.solutions[i].Cost < p.solutions[j].Cost })
	if p.limited && len(p.solutions) >= p.maxSize {
		p.solutions = p.solutions[:len(p.solutions)-1]
	}

}

// BestSolution returns a clone of the lowest-cost solution, or nil if empty.
func (p *Pool) BestSolution() *metaheuristc.RandomKeyValue {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.solutions) == 0 {
		return nil
	}

	return p.solutions[0].Clone()
}

// GetSolution returns a clone of the solution at index, panics if out of range.
func (p *Pool) GetSolution(index int) *metaheuristc.RandomKeyValue {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.solutions[index].Clone()
}

// SolutionsCount returns the number of solutions in the pool.
func (p *Pool) SolutionsCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.solutions)
}

// BestSolutionCost returns the lowest cost, or math.MaxInt if empty.
func (p *Pool) BestSolutionCost() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.solutions) == 0 {
		return math.MaxInt
	}

	return p.solutions[0].Cost
}

// Size returns the number of solutions in the pool.
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.solutions)
}
