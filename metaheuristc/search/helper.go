package search

import "fmt"

func PrintSolver(name string, s Local) {
	fmt.Printf("%s:\n", name)
	fmt.Printf("    Local Search:\n       %s\n", s.String())
}
