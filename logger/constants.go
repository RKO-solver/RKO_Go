package logger

import "strings"

// DefaultLogLevel is used when no level is explicitly configured.
const DefaultLogLevel = INFO

// Level controls how much detail a Logger prints.
type Level uint8

const (
	SILENT Level = iota
	INFO
	VERBOSE
)

// GetLogLevel parses name into a Level, defaulting to INFO when unrecognized.
func GetLogLevel(name string) Level {
	name = strings.ToUpper(name)
	switch name {
	case "SILENT":
		return SILENT
	case "INFO":
		return INFO
	case "VERBOSE":
		return VERBOSE
	default:
		return DefaultLogLevel
	}
}
