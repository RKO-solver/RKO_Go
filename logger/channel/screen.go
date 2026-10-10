package channel

import (
	"fmt"
	"strings"
)

const maxHorizontalSize = 300

const maxScreenVerticalLines = 15

const numLinesSolutionPool = 3

const freeTopLines = 1

const numMaxColumns = 3

const solverBlockSize = 100

const lastNVerboses = 3

const cleanLineCode = "\033[K"

const enterAltScreenCode = "\033[?1049h"

const leaveAltScreenCode = "\033[?1049l"

const cursorHomeCode = "\033[H"

const clearToEndCode = "\033[J"

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

// printShell redraws the dashboard in the alternate screen, entering it on the first frame.
func (d *information) printShell() {
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
	fmt.Print(b.String())
}

// printFinal leaves the alternate screen, if entered, and prints the frame once on the normal screen.
func (d *information) printFinal() {
	d.mu.Lock()
	defer d.mu.Unlock()

	var b strings.Builder
	if d.altScreen {
		d.altScreen = false
		b.WriteString(leaveAltScreenCode)
	}
	d.writeFrame(&b)
	fmt.Print(b.String())
}

// writeFrame appends the dashboard lines to b, one per line; the caller must hold d.mu.
func (d *information) writeFrame(b *strings.Builder) {
	lines := make([]string, 0)

	// Top Text
	for _ = range freeTopLines {
		lines = append(lines, cleanLineCode)
	}

	// Solution Pool Text
	idx := len(d.pool) - d.numLinesPool
	if idx < 0 {
		idx = 0
	}

	for _ = range d.numLinesPool {
		if idx < len(d.pool) {
			text := setFixedLength(fmt.Sprintf("Adding to pool %d at %.3fs", d.pool[idx].cost, d.pool[idx].time), maxHorizontalSize, false)
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
			header = setFixedLength(fmt.Sprintf("(%d) %s: %s", id, s.name, s.last.extra), solverBlockSize, false)
			info = setFixedLength(fmt.Sprintf("    Local %d Best %d", s.last.local, s.last.localBest), solverBlockSize, false)
			bottom = setFixedLength(fmt.Sprintf("    Time %.3fs", s.last.time), solverBlockSize, false)

			solversLog[id] = append(solversLog[id], header)
			solversLog[id] = append(solversLog[id], info)
			solversLog[id] = append(solversLog[id], bottom)

			numVerbose := len(d.extraMessages[id])

			if numVerbose > 0 {

				start := numVerbose - d.numVerboseMessages // lastNVerboses
				if start < 0 {
					start = 0
				}

				for i := start; i < numVerbose; i++ {
					message := setFixedLength(fmt.Sprintf("Verbose: %s Time %.3fs", d.extraMessages[id][i].message, d.extraMessages[id][i].timeStamp), solverBlockSize, false)
					solversLog[id] = append(solversLog[id], message)
				}
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
