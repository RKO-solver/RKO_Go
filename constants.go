package rko

import "strings"

// NoLimitWarning is printed when a Solver is built with neither a time limit nor an iteration limit.
const NoLimitWarning = "Warning: neither TimeLimitSeconds nor MaxIterations is set; the solver will run until its context is cancelled"

// MetaHeuristic identifies a metaheuristic.
type MetaHeuristic = int

// MetaHeuristic values; GRASP, VLNS, ALNS and IPR are not implemented yet.
const (
	MULTISTART MetaHeuristic = iota
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

// GetMetaHeuristic returns the MetaHeuristic named by label (case-insensitive), or -1 if unknown.
func GetMetaHeuristic(label string) MetaHeuristic {
	label = strings.ToUpper(label)
	switch label {
	case "MULTISTART":
		return MULTISTART
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

// GetMetaHeuristicString returns the full name of metaHeuristic, or "" if unknown.
func GetMetaHeuristicString(metaHeuristic MetaHeuristic) string {
	switch metaHeuristic {
	case MULTISTART:
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

// GetMetaHeuristicShort returns the abbreviation of metaHeuristic, or "" if unknown.
func GetMetaHeuristicShort(metaHeuristic MetaHeuristic) string {
	switch metaHeuristic {
	case MULTISTART:
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
