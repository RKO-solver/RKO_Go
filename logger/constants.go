package logger

const minimumTickerMilliseconds = 300

const defaultTickerMilliseconds = 500

const defaultBufferSize = 15

const defaultLogLevel = INFO

type Level uint8

const (
	SILENT Level = iota
	INFO
	VERBOSE
)
