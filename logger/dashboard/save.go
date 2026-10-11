package dashboard

import (
	"cmp"
	"errors"
	"fmt"
	"io"

	"github.com/RKO-solver/rko-go/logger/internal/csvio"
	"github.com/RKO-solver/rko-go/logger/internal/spill"
)

// SaveCSV writes one CSV per solver and one for the pool.
func (l *Log) SaveCSV(prefix string) error {
	prefix = cmp.Or(prefix, spill.NamePart(l.problem), "rko")

	// The lock keeps the aggregator from spilling while a solver's file and tail are read.
	l.data.mu.Lock()
	var err error
	for _, s := range l.data.solvers {
		err = errors.Join(err, csvio.Save(fmt.Sprintf("%s-%s-%d.csv", prefix, spill.NamePart(s.name), s.id),
			func(w io.Writer) error { return s.file.WriteCSV(w, s.records) }))
	}
	l.data.mu.Unlock() // Pool takes the lock itself

	pool := l.Pool()
	return errors.Join(err, csvio.Save(prefix+"-pool.csv", func(w io.Writer) error { return csvio.WritePool(w, pool) }))
}
