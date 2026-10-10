// Package solution manages pools of ranked candidate solutions.
package solution

import (
	"math"
	"slices"
	"sort"
	"sync"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/metaheuristc/rk"
	"github.com/RKO-solver/rko-go/random"
)

const defaultMaxSize = 10

// Pool is a thread-safe set of solutions sorted by ascending cost.
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

// GetGlobalInstance returns the process-wide default Pool.
func GetGlobalInstance(env definition.Environment, logger logger.Logger, rg *random.Generator) *Pool {
	once.Do(func() {
		instance = NewDefaultPool(env, rg, logger)
	})

	return instance
}

// NewPool creates a Pool of initialSize random solutions capped at maxSize.
func NewPool(maxSize int, initialSize int, env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	pool := &Pool{
		maxSize:   maxSize,
		logger:    logger,
		limited:   maxSize >= 1 && maxSize < math.MaxInt,
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

// NewDefaultPool creates a Pool of 10 random solutions capped at 10.
func NewDefaultPool(env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	return NewPool(defaultMaxSize, defaultMaxSize, env, rg, logger)
}

// NewDefaultPoolUnlimited creates an uncapped Pool of 10 random solutions.
func NewDefaultPoolUnlimited(env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	return NewPool(-1, defaultMaxSize, env, rg, logger)
}

// AddSolution adds a copy of solution if it beats the pool's best.
func (p *Pool) AddSolution(solution *metaheuristc.RandomKeyValue, time float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.solutions) == 0 {
		p.solutions = append(p.solutions, solution.Clone())
		p.logger.AddSolutionPool(solution.Cost, time)
		return
	}

	if p.solutions[0].Cost <= solution.Cost {
		return
	}

	p.solutions = slices.Insert(p.solutions, 0, solution.Clone())
	p.logger.AddSolutionPool(solution.Cost, time)
	if p.limited && len(p.solutions) >= p.maxSize {
		p.solutions = p.solutions[:p.maxSize]
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

// GetSolution returns a clone of the solution at index (0 is the best).
func (p *Pool) GetSolution(index int) *metaheuristc.RandomKeyValue {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.solutions[index].Clone()
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
