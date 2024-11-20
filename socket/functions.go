package socket

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Println("New connection established - connection from: ", ws.RemoteAddr())
	s.mu.Lock()
	s.conn[ws] = true
	s.mu.Unlock()
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
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

func (s *Server) broadcast(b []byte) {
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
