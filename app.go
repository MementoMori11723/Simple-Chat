package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	//go:embed pages/*html
	files  embed.FS
	layout = "pages/layout.html"
)

type Server struct {
	conn map[*websocket.Conn]bool
	mu   sync.Mutex
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	tmpl, err := template.ParseFS(files, layout, "pages/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error parsing template: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error executing template: ", err)
		return
	}
}

func errorHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFS(files, layout, "pages/error.html")
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("Error executing template: ", err)
		return
	}
}

func NewServer() *Server {
	return &Server{
		conn: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
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

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/error", errorHandler)
	fmt.Println("Server started at :11000")
	err := http.ListenAndServe(":11000", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
