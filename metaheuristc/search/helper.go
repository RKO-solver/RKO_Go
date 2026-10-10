package search

import "fmt"

// PrintSolver prints a solver name and its local search.
func PrintSolver(name string, s Local) {
	fmt.Printf("%s:\n", name)
	fmt.Printf("    Local Search:\n       %s\n", s.String())
}
