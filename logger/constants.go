package logger

const DefaultLogLevel = INFO

type Level uint8

const (
	SILENT Level = iota
	INFO
	VERBOSE
)
