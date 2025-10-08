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

const defaultMazSize = 200

type Pool struct {
	mu        sync.RWMutex
	solutions []*metaheuristc.RandomKeyValue
	maxSize   int
	logger    logger.Logger
}

var (
	instance *Pool
	once     sync.Once
)

func GetGlobalInstance(env definition.Environment, logger logger.Logger, rg *random.Generator) *Pool {
	once.Do(func() {
		instance = NewDefaultPool(env, rg, logger)
	})

	return instance
}

func NewPool(maxSize int, initialSize int, env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	pool := &Pool{
		maxSize:   maxSize,
		logger:    logger,
		solutions: make([]*metaheuristc.RandomKeyValue, 0, initialSize),
	}

	if initialSize > maxSize {
		maxSize = initialSize
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

func NewDefaultPool(env definition.Environment, rg *random.Generator, logger logger.Logger) *Pool {
	return NewPool(defaultMazSize, defaultMazSize, env, rg, logger)
}

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
	if len(p.solutions) >= p.maxSize {
		p.solutions = p.solutions[:len(p.solutions)-1]
	}

}

func (p *Pool) BestSolution() *metaheuristc.RandomKeyValue {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.solutions) == 0 {
		return nil
	}

	return p.solutions[0].Clone()
}

func (p *Pool) GetSolution(index int) *metaheuristc.RandomKeyValue {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.solutions[index].Clone()
}

func (p *Pool) SolutionsCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.solutions)
}

func (p *Pool) BestSolutionCost() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.solutions) == 0 {
		return math.MaxInt
	}

	return p.solutions[0].Cost
}

func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.solutions)
}
