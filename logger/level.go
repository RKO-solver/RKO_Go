package logger

import (
	"fmt"
	"strconv"
	"strings"
)

// Level controls how much detail a Logger prints.
type Level int8

const (
	LevelSilent  Level = -1 // prints nothing
	LevelInfo    Level = 0  // progress and worker start/done
	LevelVerbose Level = 1  // also events and pool improvements
)

// DefaultLevel is the level of the zero value.
const DefaultLevel = LevelInfo

// String returns the upper-case name of l.
func (l Level) String() string {
	switch l {
	case LevelSilent:
		return "SILENT"
	case LevelInfo:
		return "INFO"
	case LevelVerbose:
		return "VERBOSE"
	}
	return "Level(" + strconv.Itoa(int(l)) + ")"
}

// ParseLevel parses a level name, ignoring case.
func ParseLevel(s string) (Level, error) {
	for _, l := range [...]Level{LevelSilent, LevelInfo, LevelVerbose} {
		if strings.EqualFold(s, l.String()) {
			return l, nil
		}
	}
	return 0, fmt.Errorf("logger: unknown level %q", s)
}
