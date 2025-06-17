package logger

import (
	"errors"
	"strings"
)

func GetLevel(level string) (Level, error) {
	check := strings.ToUpper(level)
	switch check {
	case "SILENT":
		return SILENT, nil
	case "INFO":
		return INFO, nil
	case "DEBUG":
		return DEBUG, nil
	case "VERBOSE":
		return VERBOSE, nil
	case "ALL":
		return ALL, nil
	}
	return SILENT, errors.New("invalid Level")
}

func GetLevelString(level Level) string {
	switch level {
	case SILENT:
		return "Silent"
	case INFO:
		return "Info"
	case DEBUG:
		return "Debug"
	case VERBOSE:
		return "Verbose"
	case ALL:
		return "All"
	default:
		return ""
	}
}
