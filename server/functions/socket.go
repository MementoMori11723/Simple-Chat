package functions

import (
	"fmt"
	"io"
	"net/http"
	"simple-chat/server/database"
	"sync"

	"golang.org/x/net/websocket"
)

type _server struct {
	conn map[*websocket.Conn]bool
	mu   sync.Mutex
}

func socket() *_server {
	return &_server{
		conn: make(map[*websocket.Conn]bool),
	}
}

func SocketHandler() http.Handler {
  go database.Test()
  return websocket.Handler(socket().HandleWS)
}

func (s *_server) HandleWS(ws *websocket.Conn) {
  fmt.Println("route is :", ws.Request().PathValue("id"))
	fmt.Println("New connection established - connection from: ", ws.RemoteAddr())
	s.mu.Lock()
	s.conn[ws] = true
	s.mu.Unlock()
	s.readLoop(ws)
}

func (s *_server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client: ", ws.RemoteAddr())
				s.mu.Lock()
				delete(s.conn, ws)
				s.mu.Unlock()
				break
			}
			fmt.Println("Error reading from connection: ", err)
			continue
		}
		msg := buf[:n]
		s.broadcast(msg)
	}
}

func (s *_server) broadcast(b []byte) {
	s.mu.Lock()
	for ws := range s.conn {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Error broadcasting message: ", err)
			}
		}(ws)
	}
	s.mu.Unlock()
}
