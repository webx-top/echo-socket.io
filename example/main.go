package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/webx-top/echo"
	esi "github.com/webx-top/echo-socket.io"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/param"
)

//go:embed socket-client.html
var html []byte

func main() {
	e := echo.New()
	w := socketIOWrapper()
	w.Serve()
	defer w.Close()
	e.Any("/socket.io/", w)
	e.Get("/", func(c echo.Context) error {
		return c.HTML(engine.Bytes2str(html))
	})

	e.Logger().Fatal(e.Run(standard.New(":4444")))
}

func socketIOWrapper() *esi.Wrapper {
	wrapper, err := esi.NewWrapper(nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	wrapper.OnConnect("", func(context echo.Context, conn socketio.Conn) error {
		fmt.Println(`[`, time.Now().Format(param.DateTimeNormal), `]`, "connected:", conn.ID())
		return nil
	})
	wrapper.OnError("", func(context echo.Context, conn socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	wrapper.OnDisconnect("", func(context echo.Context, conn socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	wrapper.OnEvent("", "test", func(context echo.Context, conn socketio.Conn, msg string) {
		fmt.Println(`[`, time.Now().Format(param.DateTimeNormal), `]`, "notice:", msg)
		b, _ := json.MarshalIndent(context.Forms(), ``, `	`)
		fmt.Println("formData:", string(b))
		conn.Emit("test", msg) // reply message. "test" channel
	})

	return wrapper
}
