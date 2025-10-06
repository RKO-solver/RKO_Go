package logger

type messageTpe = uint8

const (
	infoMessage messageTpe = iota
	verboseMessage
)

type channelMessage struct {
	t       messageTpe
	id      int
	info    solverInfo
	message string
}
