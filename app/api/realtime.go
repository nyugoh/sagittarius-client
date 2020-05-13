package api

import (
	"fmt"
	io "github.com/googollee/go-socket.io"
)

var clients []io.Conn

func (app *App) IOServer() {
	app.SocketServer.OnConnect("/", func(s io.Conn) error {
		clients = append(clients, s)
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		s.Join("all") // Broadcast to all connected clients
		return nil
	})

	app.SocketServer.OnEvent("/", "stock-update", func(s io.Conn, msg string) {
		fmt.Println("Notice :", msg)
		s.Emit(msg)
	})

	app.SocketServer.OnEvent("/", "connection-ping", func(s io.Conn, msg string) {
		fmt.Println("Ping from:", msg)
		s.Emit(msg)
	})

	app.SocketServer.OnEvent("/", "msg", func(s io.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	app.SocketServer.OnEvent("/", "bye", func(s io.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	app.SocketServer.OnError("/", func(s io.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	app.SocketServer.OnDisconnect("/", func(s io.Conn, reason string) {
		fmt.Println("closed", reason)
	})
}

