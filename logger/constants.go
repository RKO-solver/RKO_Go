package logger

import "strings"

const DefaultLogLevel = INFO

type Level uint8

const (
	SILENT Level = iota
	INFO
	VERBOSE
)

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
