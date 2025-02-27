package functions

import (
	"io"
	"log/slog"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type _server struct {
	conn map[*websocket.Conn]bool
	mu   sync.Mutex
  maxConn int
}

func socket() *_server {
	return &_server{
		conn: make(map[*websocket.Conn]bool),
    maxConn: 3,
	}
}

func SocketHandler() http.Handler {
	return websocket.Handler(socket().HandleWS)
}

func (s *_server) HandleWS(ws *websocket.Conn) {
  if len(s.conn) >= s.maxConn {
    slog.Error("Max connections reached", "Total connections", len(s.conn))
    ws.Write([]byte("Max connections reached, try again later!"))
    ws.Close()
    return
  }
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
				s.mu.Lock()
				delete(s.conn, ws)
				s.mu.Unlock()
				break
			}
			slog.Error(
				"Error reading from connection: ",
				"Error", err.Error(),
			)
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
				slog.Error(
					"Error broadcasting message: ",
					"Error", err.Error(),
				)
			}
		}(ws)
	}
	s.mu.Unlock()
}
