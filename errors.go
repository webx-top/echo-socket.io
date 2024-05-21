package echo_socket_io

import "errors"

var (
	ErrServerCannotBeNil  = errors.New("socket.io server can not be nil")
	ErrContextCannotBeNil = errors.New("socket.io context can not be nil")
)
