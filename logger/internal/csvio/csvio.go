// Package csvio holds the CSV file helpers shared by the loggers.
package csvio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/RKO-solver/rko-go/logger"
)

// Save creates the file name, runs write on it and closes it, returning every error.
func Save(name string, write func(io.Writer) error) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	return errors.Join(write(f), f.Close())
}

// WritePool writes the pool CSV, header included, to w.
func WritePool(w io.Writer, pool []logger.PoolEntry) error {
	bw := bufio.NewWriter(w)
	bw.WriteString("cost,time\n")
	for _, e := range pool {
		fmt.Fprintf(bw, "%d,%.8f\n", e.Cost, e.Elapsed.Seconds())
	}
	return bw.Flush()
}
