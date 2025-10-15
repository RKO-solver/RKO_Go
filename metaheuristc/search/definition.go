package search

import (
	"strings"

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
	String() string
}

func GetSearchType(name string) Type {
	name = strings.ToUpper(name)
	switch name {
	case "SWAP":
		return Swap
	case "MIRROR":
		return Mirror
	case "FAREY":
		return Farey
	case "NELDER":
		return Nelder
	case "RVND":
		return RVND
	default:
		return -1
	}
}

func GetSearchString(se Type) string {
	switch se {
	case Swap:
		return "Swap"
	case Mirror:
		return "Mirror"
	case Farey:
		return "Farey"
	case Nelder:
		return "Nelder-Mead"
	case RVND:
		return "RVND"
	default:
		return ""
	}
}
