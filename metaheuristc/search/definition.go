package search

import (
	"github.com/lucasmends/rko-go/metaheuristc"
	"github.com/lucasmends/rko-go/random"
)

type Type = int

const (
	Swap Type = iota
	Mirror
	Farey
	RVND
)

type Local interface {
	SetRG(rg *random.Generator)
	Search(rko *metaheuristc.RandomKeyValue)
}
