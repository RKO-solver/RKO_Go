package dashboard

import (
	"fmt"
	"io"
	"strings"
)

func setFixedLength(s string, length int, pad bool) string {
	if len(s) > length {
		// Truncate the string
		return s[:length-3] + "..." + cleanLineCode
	}
	if pad && len(s) < length {
		// Pad the string with spaces
		padding := length - len(s)
		return s + strings.Repeat(" ", padding)
	}
	return s + cleanLineCode
}

func (d *information) printShell(w io.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var b strings.Builder
	if !d.altScreen {
		d.altScreen = true
		b.WriteString(enterAltScreenCode)
	}
	b.WriteString(cursorHomeCode)
	d.writeFrame(&b)
	b.WriteString(clearToEndCode)
	io.WriteString(w, b.String())
}

func (d *information) printFinal(w io.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var b strings.Builder
	if d.altScreen {
		d.altScreen = false
		b.WriteString(leaveAltScreenCode)
	}
	d.writeFrame(&b)
	for _, r := range d.workers {
		fmt.Fprintf(&b, "(%d) %s done cost=%d time=%.3fs\n", r.ID, r.Name, r.Cost, r.Elapsed.Seconds())
	}
	io.WriteString(w, b.String())
}

func (d *information) writeFrame(b *strings.Builder) {
	lines := make([]string, 0)

	// Top Text
	for range freeTopLines {
		lines = append(lines, cleanLineCode)
	}

	// Solution Pool Text
	idx := len(d.pool) - d.numLinesPool
	if idx < 0 {
		idx = 0
	}

	for range d.numLinesPool {
		if idx < len(d.pool) {
			text := setFixedLength(fmt.Sprintf("pool cost=%d time=%.3fs", d.pool[idx].Cost, d.pool[idx].Elapsed.Seconds()), maxHorizontalSize, false)
			lines = append(lines, text)
			idx++
		} else {
			break
		}
	}

	lines = append(lines, cleanLineCode)

	var header, info, bottom string

	solversLog := make([][]string, 0)

	for id, s := range d.solvers {
		solversLog = append(solversLog, []string{})
		if s.hasLast { // read last, not records: records may have just been spilled and emptied
			header = setFixedLength(fmt.Sprintf("(%d) %s: iteration %d", id, s.name, s.last.iteration), solverBlockSize, false)
			info = setFixedLength(fmt.Sprintf("    Local %d Best %d", s.last.local, s.last.best), solverBlockSize, false)
			bottom = setFixedLength(fmt.Sprintf("    Time %.3fs", s.last.elapsed.Seconds()), solverBlockSize, false)

			solversLog[id] = append(solversLog[id], header)
			solversLog[id] = append(solversLog[id], info)
			solversLog[id] = append(solversLog[id], bottom)

			for _, e := range s.events {
				message := setFixedLength(fmt.Sprintf("    %s time=%.3fs", e.attr, e.elapsed.Seconds()), solverBlockSize, false)
				solversLog[id] = append(solversLog[id], message)
			}
		}
	}

	for _, log := range solversLog {
		for _, line := range log {
			lines = append(lines, line)
		}
		lines = append(lines, cleanLineCode)
	}

	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
}
