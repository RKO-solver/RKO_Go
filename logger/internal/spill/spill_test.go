package spill

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/RKO-solver/rko-go/logger"
	"github.com/stretchr/testify/require"
)

// Checks Capacity, that All and WriteCSV return file records then memory in order, and that Remove deletes the file.
func TestWriteThenAllAndWriteCSV(t *testing.T) {
	require.Equal(t, 0, Capacity(0))
	require.Equal(t, 0, Capacity(-1))
	require.Equal(t, 1, Capacity(10))
	require.Equal(t, 5461, Capacity(128*1024))

	dir := t.TempDir()
	var h File
	h.Reset(Config{Dir: dir, Bytes: 3 * recordSize})
	t.Cleanup(func() { require.NoError(t, h.Remove()) })

	first := []logger.Sample{{Local: -3, Best: -5, Elapsed: 1250 * time.Millisecond}, {Local: 4, Best: 2, Elapsed: 2500 * time.Millisecond}}
	second := []logger.Sample{{Local: 7, Best: 2, Elapsed: 3 * time.Second}}
	memory := []logger.Sample{{Local: 1, Best: 1, Elapsed: 4125 * time.Millisecond}}
	require.True(t, h.Write(first))
	require.True(t, h.Write(second))
	require.Equal(t, slices.Concat(first, second, memory), h.All(memory))

	var out bytes.Buffer
	require.NoError(t, h.WriteCSV(&out, memory))
	require.Equal(t, "best,local,time\n-5,-3,1.25000000\n2,4,2.50000000\n2,7,3.00000000\n1,1,4.12500000\n", out.String())

	name := h.f.Name()
	require.NoError(t, h.Remove())
	_, err := os.Stat(name)
	require.ErrorIs(t, err, fs.ErrNotExist)
}

// Checks that a failed create turns spilling off once and All still returns the in-memory records.
func TestCreateErrorKeepsHistoryInMemory(t *testing.T) {
	var h File
	h.Reset(Config{Dir: filepath.Join(t.TempDir(), "missing"), Bytes: 3 * recordSize})
	recs := []logger.Sample{{Local: 1, Best: 1, Elapsed: 500 * time.Millisecond}}

	require.False(t, h.Write(recs)) // prints the one warning
	require.True(t, h.off)
	require.False(t, h.Write(recs)) // off: no second attempt, no second warning
	require.Nil(t, h.f)
	require.Equal(t, recs, h.All(recs))
}

// Checks that a failed write turns spilling off and keeps the count of records already in the file.
func TestWriteErrorKeepsWrittenRecords(t *testing.T) {
	var h File
	h.Reset(Config{Dir: t.TempDir(), Bytes: 3 * recordSize})
	t.Cleanup(func() { _ = h.Remove() }) // f is already closed: only the delete matters

	require.True(t, h.Write([]logger.Sample{{Local: 1}}))
	require.NoError(t, h.f.Close())
	require.False(t, h.Write([]logger.Sample{{Local: 2}}))
	require.True(t, h.off)
	require.Equal(t, 1, h.records)
}

// Checks the file name pattern and the skipping of empty parts.
func TestPatternSanitizesNameParts(t *testing.T) {
	require.Equal(t, "rko-history-tsp-BRKGA-2-*.bin", pattern(Config{Problem: "tsp", Solver: "BRKGA", ID: 2}))
	require.Equal(t, "rko-history-BRKGA-0-*.bin", pattern(Config{Solver: "BRKGA"}))
	require.Equal(t, "rko-history-0-*.bin", pattern(Config{}))
}

// Checks that the first Write creates the file in the configured directory with the sanitized, truncated name.
func TestWriteCreatesNamedFile(t *testing.T) {
	dir := t.TempDir()
	var h File
	h.Reset(Config{Dir: dir, Bytes: recordSize, Problem: "inst/a:b" + strings.Repeat("x", 100), Solver: "BRKGA", ID: 3})
	t.Cleanup(func() { require.NoError(t, h.Remove()) })

	require.True(t, h.Write([]logger.Sample{{Local: 1}}))
	require.Equal(t, dir, filepath.Dir(h.f.Name()))
	base := filepath.Base(h.f.Name())
	want := "rko-history-" + ("inst_a_b" + strings.Repeat("x", 100))[:maxNamePart] + "-BRKGA-3-"
	require.True(t, strings.HasPrefix(base, want), base)
	require.True(t, strings.HasSuffix(base, ".bin"), base)
}

// Checks that Write after Remove creates no file until the next Reset.
func TestRemoveTurnsSpillingOffUntilReset(t *testing.T) {
	dir := t.TempDir()
	var h File
	h.Reset(Config{Dir: dir, Bytes: recordSize})
	t.Cleanup(func() { require.NoError(t, h.Remove()) })
	recs := []logger.Sample{{Local: 1}}

	require.NoError(t, h.Remove()) // before any file
	require.False(t, h.Write(recs))
	require.Nil(t, h.f)

	h.Reset(Config{Dir: dir, Bytes: recordSize})
	require.True(t, h.Write(recs))
	require.NoError(t, h.Remove())
	require.False(t, h.Write(recs))
	require.Nil(t, h.f)
}
