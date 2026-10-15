// Package policyio prints and saves the tuners' policies recorded by the loggers.
package policyio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/RKO-solver/rko-go/logger/internal/csvio"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// Metaheuristic is the policy recorded by one metaheuristic.
type Metaheuristic struct {
	ID    int
	Name  string
	Table logger.Table
}

// Print writes each metaheuristic's policy to w under a "(id) name policy" heading.
func Print(w io.Writer, metaheuristics []Metaheuristic) error {
	// bufio.Writer keeps its first error and returns it from Flush, so per-call errors are not checked.
	bw := bufio.NewWriter(w)
	for i, s := range metaheuristics {
		if i > 0 {
			bw.WriteByte('\n')
		}
		fmt.Fprintf(bw, "(%d) %s policy\n", s.ID, s.Name)
		tw := tabwriter.NewWriter(bw, 0, 0, 2, ' ', 0) // columns two spaces apart
		writeRow(tw, s.Table.Header)
		for _, row := range s.Table.Rows {
			writeRow(tw, row)
		}
		tw.Flush()
	}
	return bw.Flush()
}

func writeRow(w io.Writer, cells []string) {
	io.WriteString(w, strings.Join(cells, "\t")+"\n")
}

// Save writes each metaheuristic's policy to <prefix>-<metaheuristic>-<id>-policy.csv and returns every error joined.
func Save(prefix string, metaheuristics []Metaheuristic) error {
	var err error
	for _, s := range metaheuristics {
		err = errors.Join(err, csvio.Save(fmt.Sprintf("%s-%s-%d-policy.csv", prefix, spill.NamePart(s.Name), s.ID),
			func(w io.Writer) error { return csvio.WriteTable(w, s.Table) }))
	}
	return err
}
