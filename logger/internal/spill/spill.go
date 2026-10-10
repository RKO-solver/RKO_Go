// Package spill appends a solver's history records to a temporary binary file.
// The loggers fill a fixed-capacity slice and hand it to Write when full; the file
// is created on the first Write and kept open. All and SaveCSV read the file back,
// then the in-memory tail. Reset deletes the previous run's file; Remove deletes the current one
// and turns spilling off until the next Reset.
package spill

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/RKO-solver/rko-go/logger"
)

// recordSize is the size of one record: local int64, best int64, time float64, little-endian.
const recordSize = 24

// maxNamePart caps the bytes kept from each file name part.
const maxNamePart = 32

// Config sets a solver's spill directory, buffer size and file name parts.
type Config struct {
	Dir     string // "" means the system temp directory
	Bytes   int    // 0 or less keeps the history in memory
	Problem string
	Solver  string
	ID      int
}

// Capacity returns how many records fit in bytes, at least 1 when bytes is positive.
func Capacity(bytes int) int {
	if bytes <= 0 {
		return 0
	}
	return max(bytes/recordSize, 1)
}

// File holds one solver's spilled records; Reset sets it up before the first Write.
// It has no lock of its own: the caller serializes Write, reads, Reset and Remove.
type File struct {
	dir     string   // where the file is created; "" means the system temp directory
	pattern string   // os.CreateTemp pattern, built by Reset
	off     bool     // spilling disabled until the next Reset: buffer size 0 or less, a failed create or write, or Remove
	f       *os.File // nil until the first Write
	records int      // records fully written to f
	buf     []byte   // reusable encoding buffer
}

// Reset deletes the file, if any, and applies c.
func (h *File) Reset(c Config) {
	_ = h.Remove() // a file that cannot be deleted is left behind; the new run starts clean either way
	h.dir = c.Dir
	h.pattern = pattern(c)
	h.off = c.Bytes <= 0
}

// pattern returns the os.CreateTemp pattern rko-history-<problem>-<solver>-<id>-*.bin, leaving out empty name parts.
func pattern(c Config) string {
	var b strings.Builder
	b.WriteString("rko-history-")
	for _, s := range [...]string{c.Problem, c.Solver} {
		if s = namePart(s); s != "" {
			b.WriteString(s)
			b.WriteByte('-')
		}
	}
	b.WriteString(strconv.Itoa(c.ID))
	b.WriteString("-*.bin")
	return b.String()
}

// namePart replaces path separators, characters some file systems reject and control characters with '_', and cuts s to maxNamePart bytes.
func namePart(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || strings.ContainsRune(`/\:*?"<>|`, r) {
			return '_'
		}
		return r
	}, s)
	if len(s) > maxNamePart {
		n := maxNamePart
		for n > 0 && !utf8.RuneStart(s[n]) {
			n--
		}
		s = s[:n]
	}
	return s
}

// Write appends recs to the file in one write, creating the file first if needed, and reports success.
// On false the caller keeps recs in memory; after a failed create or write, or Remove, Write returns false until the next Reset.
func (h *File) Write(recs []logger.Data) bool {
	if h.off {
		return false
	}
	if h.f == nil {
		f, err := os.CreateTemp(h.dir, h.pattern)
		if err != nil {
			h.fail(err)
			return false
		}
		h.f = f
	}
	h.buf = slices.Grow(h.buf[:0], len(recs)*recordSize) // reuse the encoding buffer across flushes
	for _, r := range recs {
		// The uint64 conversions keep the two's-complement bits, so negative costs round-trip.
		h.buf = binary.LittleEndian.AppendUint64(h.buf, uint64(r.LocalCost))
		h.buf = binary.LittleEndian.AppendUint64(h.buf, uint64(r.BestCost))
		h.buf = binary.LittleEndian.AppendUint64(h.buf, math.Float64bits(r.Time))
	}
	if _, err := h.f.Write(h.buf); err != nil {
		h.fail(err) // bytes of a partial write lie past records and are never read
		return false
	}
	h.records += len(recs)
	return true
}

// fail turns spilling off and prints one warning; an open file is kept so its records can still be read back.
func (h *File) fail(err error) {
	h.off = true
	fmt.Fprintf(os.Stderr, "Warning: solver history stays in memory: %v\n", err)
}

// forEach calls fn for every record in the file, in write order.
func (h *File) forEach(fn func(logger.Data)) error {
	if h.records == 0 {
		return nil
	}
	// SectionReader uses ReadAt, so the append offset of h.f is untouched and Write can continue after a read.
	r := bufio.NewReaderSize(io.NewSectionReader(h.f, 0, int64(h.records)*recordSize), 64<<10)
	var rec [recordSize]byte
	for range h.records {
		if _, err := io.ReadFull(r, rec[:]); err != nil {
			return err
		}
		fn(logger.Data{
			LocalCost: int(int64(binary.LittleEndian.Uint64(rec[0:8]))),
			BestCost:  int(int64(binary.LittleEndian.Uint64(rec[8:16]))),
			Time:      math.Float64frombits(binary.LittleEndian.Uint64(rec[16:24])),
		})
	}
	return nil
}

// All returns the file's records followed by memory, in a new slice.
func (h *File) All(memory []logger.Data) []logger.Data {
	all := make([]logger.Data, 0, h.records+len(memory))
	if err := h.forEach(func(d logger.Data) { all = append(all, d) }); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cannot read solver history file: %v\n", err)
	}
	return append(all, memory...)
}

// SaveCSV streams the header best,local,time, the file's records and then memory to filename.
func (h *File) SaveCSV(filename string, memory []logger.Data) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	// bufio.Writer keeps its first error and returns it from Flush, so per-call errors are not checked.
	w := bufio.NewWriterSize(out, 64<<10)
	w.WriteString("best,local,time\n")
	var line []byte
	row := func(d logger.Data) {
		line = strconv.AppendInt(line[:0], int64(d.BestCost), 10)
		line = append(line, ',')
		line = strconv.AppendInt(line, int64(d.LocalCost), 10)
		line = append(line, ',')
		line = strconv.AppendFloat(line, d.Time, 'f', 8, 64)
		line = append(line, '\n')
		w.Write(line)
	}
	err = h.forEach(row)
	for _, d := range memory {
		row(d)
	}
	return errors.Join(err, w.Flush(), out.Close())
}

// Remove closes and deletes the file, if any, and turns spilling off until the next Reset.
func (h *File) Remove() error {
	h.off = true
	if h.f == nil {
		return nil
	}
	err := errors.Join(h.f.Close(), os.Remove(h.f.Name()))
	h.f = nil
	h.records = 0
	return err
}
