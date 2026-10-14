package dashboard

import "time"

const minimumRefresh = 300 * time.Millisecond

const defaultRefresh = 500 * time.Millisecond

const defaultBufferSize = 1024

const defaultPoolLines = 3

// defaultHistoryBytes is each metaheuristic's default history buffer, in bytes, before it spills to a file.
const defaultHistoryBytes = 128 * 1024

const maxHorizontalSize = 300

const freeTopLines = 1

const metaheuristicBlockSize = 100

const cleanLineCode = "\033[K"

const enterAltScreenCode = "\033[?1049h"

const leaveAltScreenCode = "\033[?1049l"

const cursorHomeCode = "\033[H"

const clearToEndCode = "\033[J"
