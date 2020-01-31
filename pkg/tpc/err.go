package tpc

import "errors"

var (
	ErrAlreadyConnected    = errors.New("already connected")
	ErrAlreadyDisconnected = errors.New("already disconnected")
)
