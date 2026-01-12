package rko

import "strings"

type MetaHeuristic = int

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

func GetMetaHeuristicString(metaHeuristic MetaHeuristic) string {
	switch metaHeuristic {
	case MULTISTART:
		return "MultiStart"
	case SA:
		return "Simulation Annealing"
	case GA:
		return "Genetic Algorithm"
	case VNS:
		return "Variable Neighborhood Search"
	case ILS:
		return "Iterated Search"
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
