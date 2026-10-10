package channel

const minimumTickerMilliseconds = 300

const defaultTickerMilliseconds = 500

const defaultBufferSize = 15

// defaultHistoryBytes is each solver's default history buffer, in bytes, before it spills to a file.
const defaultHistoryBytes = 128 * 1024
