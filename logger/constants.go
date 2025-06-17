package logger

type Level uint8

const (
	SILENT Level = iota
	INFO
	DEBUG
	VERBOSE
	ALL
)
