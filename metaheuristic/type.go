package metaheuristic

import "strings"

// Type identifies a metaheuristic.
type Type int

// Type values; GRASP, VLNS, ALNS and IPR are not implemented yet.
const (
	MultiStart Type = iota
	SA
	GA
	VNS
	ILS
	BRKGA
	GRASP
	LNS
	VLNS
	ALNS
	IPR
)

// ParseType returns the Type named by label (case-insensitive), or -1 if unknown.
func ParseType(label string) Type {
	label = strings.ToUpper(label)
	switch label {
	case "MULTISTART":
		return MultiStart
	case "SA":
		return SA
	case "GA":
		return GA
	case "VNS":
		return VNS
	case "ILS":
		return ILS
	case "BRKGA":
		return BRKGA
	case "GRASP":
		return GRASP
	case "LNS":
		return LNS
	case "VLNS":
		return VLNS
	case "ALNS":
		return ALNS
	case "IPR":
		return IPR
	default:
		return -1
	}
}

// String returns the full name of t, or "" if unknown.
func (t Type) String() string {
	switch t {
	case MultiStart:
		return "MultiStart"
	case SA:
		return "Simulated Annealing"
	case GA:
		return "Genetic Algorithm"
	case VNS:
		return "Variable Neighbourhood Search"
	case ILS:
		return "Iterated Local Search"
	case BRKGA:
		return "Biased Random Key Genetic Algorithm"
	case GRASP:
		return "Greedy Randomized Adaptive Search Procedure"
	case LNS:
		return "Large Neighbourhood Search"
	case VLNS:
		return "Very Large Neighbourhood Search"
	case ALNS:
		return "Adaptive Large Neighbourhood Search"
	case IPR:
		return "Implicit Path Relinking"
	default:
		return ""
	}
}

// Short returns the abbreviation of t, or "" if unknown.
func (t Type) Short() string {
	switch t {
	case MultiStart:
		return "MS"
	case SA:
		return "SA"
	case GA:
		return "GA"
	case VNS:
		return "VNS"
	case ILS:
		return "ILS"
	case BRKGA:
		return "BRKGA"
	case GRASP:
		return "GRASP"
	case LNS:
		return "LNS"
	case VLNS:
		return "VLNS"
	case ALNS:
		return "ALNS"
	case IPR:
		return "IPR"
	default:
		return ""
	}
}

// Implemented reports whether New can build t.
func (t Type) Implemented() bool {
	switch t {
	case MultiStart, SA, GA, VNS, ILS, BRKGA, LNS:
		return true
	}
	return false
}
