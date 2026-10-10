// Package search implements the local searches used by the metaheuristics.
package search

import (
	"context"
	"strings"

	"github.com/RKO-solver/rko-go/metaheuristc"
	"github.com/RKO-solver/rko-go/random"
)

// Type identifies a local search.
type Type = int

const (
	// Swap swaps keys within each Environment.SwapSearch range.
	Swap Type = iota
	// Mirror replaces a key k by 1-k.
	Mirror
	// Farey tries a value from each Farey interval for each key.
	Farey
	// Nelder runs Nelder-Mead over the solution and two pool solutions.
	Nelder
	// RVND is a random variable neighbourhood descent over other searches.
	RVND
)

// Local is a local search that improves a solution in place.
type Local interface {
	// SetRG sets the random generator of the search.
	SetRG(rg *random.Generator)
	// Search improves rko in place, stopping early when ctx is done.
	Search(ctx context.Context, rko *metaheuristc.RandomKeyValue)
	// String returns the search name.
	String() string
}

// GetSearchType returns the Type named by name (case-insensitive), or -1 if unknown.
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

// GetSearchString returns the display name of se, or "" if unknown.
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
