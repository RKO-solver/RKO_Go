package search

import (
	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

type Type = int

const (
	Swap Type = iota
	Mirror
	Farey
	Nelder
	RVND
)

type Local interface {
	SetRG(rg *random.Generator)
	Search(rko *metaheuristc.RandomKeyValue)
}
