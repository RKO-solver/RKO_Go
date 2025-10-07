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
	case "VERBOSE":
		return VERBOSE, nil
	}
	return SILENT, errors.New("invalid Level")
}

func GetLevelString(level Level) string {
	switch level {
	case SILENT:
		return "Silent"
	case INFO:
		return "Info"
	case VERBOSE:
		return "Verbose"
	default:
		return ""
	}
}
