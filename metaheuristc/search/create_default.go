package search

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/random"
)

func Create(typeSearch Type, environment definition.Environment, rg *random.Generator) Local {
	switch typeSearch {
	case Swap:
		return CreateSwapLocalSearch(environment)
	case Mirror:
		return CreateMirrorLocalSearch(environment)
	case Farey:
		return CreateFareyLocalSearch(environment, rg)
	case RVND:
		neighbourhood := make([]Type, 3)
		neighbourhood[0] = Swap
		neighbourhood[1] = Mirror
		neighbourhood[2] = Farey
		return CreateRVND(environment, rg, neighbourhood)
	default:
		return CreateDefault(environment, rg)
	}
}

func CreateDefault(environment definition.Environment, rg *random.Generator) Local {
	neighbourhood := make([]Type, 3)
	neighbourhood[0] = Swap
	neighbourhood[1] = Mirror
	neighbourhood[2] = Farey

	return CreateRVND(environment, rg, neighbourhood)
}
