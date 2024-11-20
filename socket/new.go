package socket 

import (
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	conn map[*websocket.Conn]bool
	mu   sync.Mutex
}

func NewServer() *Server {
	return &Server{
		conn: make(map[*websocket.Conn]bool),
	}
}
