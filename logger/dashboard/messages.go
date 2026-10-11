package dashboard

import (
	"log/slog"
	"time"
)

type messageType uint8

const (
	infoMessage messageType = iota
	eventMessage
)

// solverInfo is one Register call sent through the update channel.
type solverInfo struct {
	local, best, iteration int
	elapsed                time.Duration
}

// event is one Event call sent through the update channel.
type event struct {
	attr    slog.Attr
	elapsed time.Duration
}

type message struct {
	t     messageType
	id    int
	info  solverInfo
	event event
}
