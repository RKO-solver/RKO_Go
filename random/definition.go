package random

import (
	"math/rand/v2"
	"sync"
	"time"
)

type Generator struct {
	rand *rand.Rand
}

var (
	instance *Generator
	once     sync.Once
)

func NewGenerator() *Generator {
	currentTimeSeed := uint64(time.Now().UnixNano())
	source := rand.NewPCG(currentTimeSeed, currentTimeSeed+1)

	return &Generator{
		rand: rand.New(source),
	}
}

func NewGeneratorWithSeed(seed uint64) *Generator {
	source := rand.NewPCG(seed, seed+1)

	return &Generator{
		rand: rand.New(source),
	}
}

func GetGlobalInstance() *Generator {
	once.Do(func() {
		instance = NewGenerator()
	})

	return instance
}

func NewGeneratorSeed(seed uint64) *Generator {
	source := rand.NewPCG(seed, seed+1)

	return &Generator{
		rand: rand.New(source),
	}
}

func (g *Generator) Float64() float64 {
	return g.rand.Float64()
}

func (g *Generator) Float32() float32 {
	return g.rand.Float32()
}

func (g *Generator) IntN(n int) int {
	return g.rand.IntN(n)
}

func (g *Generator) RangeInts(maxInt, numElem int) []int {
	values := make([]int, numElem)
	for i := 0; i < numElem; i++ {
		values[i] = g.rand.IntN(maxInt)
	}
	return values
}

func (g *Generator) Permutation(n int) []int {
	return g.rand.Perm(n)
}

func (g *Generator) RangeFloat64(min, max float64) float64 {
	return min + g.rand.Float64()*(max-min)
}

func (g *Generator) RangeInt(min, max int) int {
	return min + g.rand.IntN(max-min)
}
