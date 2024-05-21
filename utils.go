package echo_socket_io

import (
	"context"

	socketio "github.com/googollee/go-socket.io"
	"github.com/webx-top/echo"
)

func getContext(conn socketio.Conn) echo.Context {
	ctx := conn.Context()
	if ctx == nil {
		return nil
	}
	switch v := ctx.(type) {
	case echo.Context:
		return v
	case context.Context:
		eCtx, ok := echo.FromStdContext(v)
		if ok {
			return eCtx
		}
	}
	return nil
}

func getContextByStd(ctx context.Context) echo.Context {
	switch v := ctx.(type) {
	case echo.Context:
		return v
	case context.Context:
		eCtx, ok := echo.FromStdContext(v)
		if ok {
			return eCtx
		}
	}
	return nil
}
