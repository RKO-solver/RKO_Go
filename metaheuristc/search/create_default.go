package search

import (
	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristc/solution"
	"github.com/RKO-solver/rko-go/random"
)

func Create(typeSearch Type, environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local {
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
		return CreateRVND(environment, solutionPool, rg, neighbourhood)
	case Nelder:
		return CreateNelderMeadLocalSearch(environment, solutionPool, rg)
	default:
		return CreateDefault(environment, solutionPool, rg)
	}
}

func CreateDefault(environment definition.Environment, solutionPool *solution.Pool, rg *random.Generator) Local {
	neighbourhood := make([]Type, 3)
	neighbourhood[0] = Swap
	neighbourhood[1] = Mirror
	neighbourhood[2] = Farey
	neighbourhood[3] = Nelder

	return CreateRVND(environment, solutionPool, rg, neighbourhood)
}
