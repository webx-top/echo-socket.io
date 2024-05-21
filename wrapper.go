package echo_socket_io

import (
	"net/http"

	"github.com/admpub/log"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/webx-top/echo"
)

// Socket.io wrapper interface
type IWrapper interface {
	OnConnect(nsp string, f func(echo.Context, socketio.Conn) error)
	OnDisconnect(nsp string, f func(echo.Context, socketio.Conn, string))
	OnError(nsp string, f func(echo.Context, socketio.Conn, error))
	OnEvent(nsp, event string, f func(echo.Context, socketio.Conn, string))
	Serve()
	Close()
	echo.Handler
}

var _ IWrapper = (*Wrapper)(nil)

func ConnInitor(r *http.Request, c engineio.Conn) {
	c.SetContext(getContextByStd(r.Context()))
}

type Wrapper struct {
	Server *socketio.Server
}

// Create wrapper and Socket.io server
func NewWrapper(options *engineio.Options) *Wrapper {
	if options == nil {
		options = &engineio.Options{
			ConnInitor: ConnInitor,
		}
	} else {
		if options.ConnInitor == nil {
			options.ConnInitor = ConnInitor
		} else {
			customInitor := options.ConnInitor
			options.ConnInitor = func(r *http.Request, c engineio.Conn) {
				ConnInitor(r, c)
				customInitor(r, c)
			}
		}
	}
	server := socketio.NewServer(options)

	return &Wrapper{
		Server: server,
	}
}

// Create wrapper with exists Socket.io server
func NewWrapperWithServer(server *socketio.Server) (*Wrapper, error) {
	if server == nil {
		return nil, ErrServerCannotBeNil
	}

	return &Wrapper{
		Server: server,
	}, nil
}

// On Socket.io client connect
func (s *Wrapper) OnConnect(nsp string, f func(echo.Context, socketio.Conn) error) {
	s.Server.OnConnect(nsp, func(conn socketio.Conn) error {
		ctx := getContext(conn)
		if ctx == nil {
			return ErrContextCannotBeNil
		}
		return f(ctx, conn)
	})
}

// On Socket.io client disconnect
func (s *Wrapper) OnDisconnect(nsp string, f func(echo.Context, socketio.Conn, string)) {
	s.Server.OnDisconnect(nsp, func(conn socketio.Conn, msg string) {
		ctx := getContext(conn)
		if ctx == nil {
			return
		}
		f(ctx, conn, msg)
	})
}

// On Socket.io error
func (s *Wrapper) OnError(nsp string, f func(echo.Context, socketio.Conn, error)) {
	s.Server.OnError(nsp, func(conn socketio.Conn, err error) {
		ctx := getContext(conn)
		if ctx == nil {
			return
		}
		f(ctx, conn, err)
	})
}

// On Socket.io event from client
func (s *Wrapper) OnEvent(nsp, event string, f func(echo.Context, socketio.Conn, string)) {
	s.Server.OnEvent(nsp, event, func(conn socketio.Conn, msg string) {
		ctx := getContext(conn)
		if ctx == nil {
			return
		}
		f(ctx, conn, msg)
	})
}

// On Socket.io event from client
func (s *Wrapper) OnEventAndReturn(nsp, event string, f func(echo.Context, socketio.Conn, string) string) {
	s.Server.OnEvent(nsp, event, func(conn socketio.Conn, msg string) string {
		ctx := getContext(conn)
		if ctx == nil {
			return ErrContextCannotBeNil.Error()
		}
		return f(ctx, conn, msg)
	})
}

// Handler function
func (s *Wrapper) Handle(context echo.Context) error {
	s.Server.ServeHTTP(context.Response().StdResponseWriter(), context.Request().StdRequest().WithContext(echo.AsStdContext(context)))
	return nil
}

func (s *Wrapper) Serve() {
	go func() {
		err := s.Server.Serve()
		if err != nil {
			log.Errorf(`socket.io: %s`, err.Error())
		}
	}()
}

func (s *Wrapper) Close() {
	err := s.Server.Close()
	if err != nil {
		log.Errorf(`socket.io: %s`, err.Error())
	}
}
