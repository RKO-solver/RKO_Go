package rko

import "strings"

// NoLimitWarning is printed when a Solver is built with neither a time limit nor an iteration limit.
const NoLimitWarning = "Warning: neither a time limit nor an iteration limit is set; the solver will run until its context is cancelled"

// Metaheuristic identifies a metaheuristic.
type Metaheuristic int

// Metaheuristic values; GRASP, VLNS, ALNS and IPR are not implemented yet.
const (
	MultiStart Metaheuristic = iota
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

// ParseMetaheuristic returns the Metaheuristic named by label (case-insensitive), or -1 if unknown.
func ParseMetaheuristic(label string) Metaheuristic {
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

// String returns the full name of m, or "" if unknown.
func (m Metaheuristic) String() string {
	switch m {
	case MultiStart:
		return "MultiStart"
	case SA:
		return "Simulated Annealing"
	case GA:
		return "Genetic Algorithm"
	case VNS:
		return "Variable Neighborhood Search"
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

// Short returns the abbreviation of m, or "" if unknown.
func (m Metaheuristic) Short() string {
	switch m {
	case MultiStart:
		return "MH"
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
