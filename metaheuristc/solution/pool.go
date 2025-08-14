package solution

import (
	"fmt"
	"github.com/lucasmends/rko-go/logger"
	"github.com/lucasmends/rko-go/logger/basic"
	"github.com/lucasmends/rko-go/metaheuristc"
	"math"
	"sort"
	"sync"
)

const defaultMazSize = 200

type Pool struct {
	mu        sync.RWMutex
	solutions []*metaheuristc.RandomKeyValue
	maxSize   int
	logger    *logger.Log
}

var (
	instance *Pool
	once     sync.Once
)

func GetGlobalInstance() *Pool {
	once.Do(func() {
		lo := logger.CreateLogger(logger.INFO, false, basic.CreateLogger())
		instance = NewDefaultPool(lo)
	})

	return instance
}

func NewPool(maxSize int, logger *logger.Log) *Pool {
	return &Pool{
		maxSize:   maxSize,
		logger:    logger,
		solutions: make([]*metaheuristc.RandomKeyValue, 0),
	}
}

func NewDefaultPool(logger *logger.Log) *Pool {
	return NewPool(defaultMazSize, logger)
}

func (p *Pool) AddSolution(solution *metaheuristc.RandomKeyValue) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.solutions) == 0 {
		p.solutions = append(p.solutions, solution)
		p.logger.Info(fmt.Sprintf("Adding solution cost %d to the pool", solution.Cost))
		return
	}

	// np better than the worst ignore
	if p.solutions[len(p.solutions)-1].Cost < solution.Cost {
		return
	}

	p.logger.Info(fmt.Sprintf("Adding solution cost %d to the pool", solution.Cost))
	p.append(solution)

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

func (p *Pool) append(solution *metaheuristc.RandomKeyValue) {
	p.solutions = append(p.solutions, solution)
	sort.Slice(p.solutions, func(i, j int) bool { return p.solutions[i].Cost < p.solutions[j].Cost })
}
