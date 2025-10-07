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

func (d *information) printShell() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.init {
		d.init = true
	} else {
		fmt.Print(fmt.Sprintf("\033[%dA", d.previousLineCount))
	}

	lines := make([]string, 0)

	// Top Text
	for _ = range freeTopLines {
		lines = append(lines, cleanLineCode)
	}

	// Solution Pool Text
	idx := len(d.pool) - numLinesSolutionPool
	if idx < 0 {
		idx = 0
	}

	for _ = range numLinesSolutionPool {
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

	for id, solver := range d.solvers {
		solversLog = append(solversLog, []string{})
		lastIdMessage := len(solver) - 1
		if lastIdMessage >= 0 {
			header = setFixedLength(fmt.Sprintf("(%d) %s: %s", solver[lastIdMessage].id, solver[lastIdMessage].name, solver[lastIdMessage].extra), solverBlockSize, false)
			info = setFixedLength(fmt.Sprintf("    Local %d Best %d", solver[lastIdMessage].local, solver[lastIdMessage].localBest), solverBlockSize, false)
			bottom = setFixedLength(fmt.Sprintf("    Time %.3fs", solver[lastIdMessage].time), solverBlockSize, false)

			solversLog[id] = append(solversLog[id], header)
			solversLog[id] = append(solversLog[id], info)
			solversLog[id] = append(solversLog[id], bottom)

			numVerbose := len(d.extraMessages[id])

			if numVerbose > 0 {

				start := numVerbose - lastNVerboses
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
		fmt.Println(line)
	}

	d.previousLineCount = len(lines)

}
