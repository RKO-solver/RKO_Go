// Package search implements the local searches used by the metaheuristics.
package search

import (
	"strings"
)

// Type identifies a local search.
type Type int

const (
	// Swap swaps keys within each Environment.SwapSearch range.
	Swap Type = iota
	// Mirror replaces a key k by 1-k.
	Mirror
	// Farey tries a value from each Farey interval for each key.
	Farey
	// NelderMead runs Nelder-Mead over the solution and two pool solutions.
	NelderMead
	// RVND is a random variable neighbourhood descent over other searches.
	RVND
)

// ParseType returns the Type named by name (case-insensitive), or -1 if unknown.
func ParseType(name string) Type {
	name = strings.ToUpper(name)
	switch name {
	case "SWAP":
		return Swap
	case "MIRROR":
		return Mirror
	case "FAREY":
		return Farey
	case "NELDER":
		return NelderMead
	case "RVND":
		return RVND
	default:
		return -1
	}
}

// String returns the display name of t, or "" if unknown.
func (t Type) String() string {
	switch t {
	case Swap:
		return "Swap"
	case Mirror:
		return "Mirror"
	case Farey:
		return "Farey"
	case NelderMead:
		return "Nelder-Mead"
	case RVND:
		return "RVND"
	default:
		return ""
	}
}
