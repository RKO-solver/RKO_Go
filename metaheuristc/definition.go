package metaheuristc

import (
	"github.com/lucasmends/rko-go/definition"
	"sort"
)

type RandomKeyValue struct {
	RK   definition.RandomKey
	Cost int
}

func (rko *RandomKeyValue) Clone() *RandomKeyValue {
	sol := &RandomKeyValue{
		RK:   make(definition.RandomKey, len(rko.RK)),
		Cost: rko.Cost,
	}
	copy(sol.RK, rko.RK)

	return sol
}

func Sort(rkos []*RandomKeyValue) {
	sort.Slice(rkos, func(i, j int) bool { return rkos[i].Cost < rkos[j].Cost })
}
