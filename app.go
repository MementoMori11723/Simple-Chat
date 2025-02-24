package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"simple-chat/server"
)

func main() {
  server := http.Server{
    Addr:    ":11000",
    Handler: Server.New(),
  }
  slog.Info(
    "Server started on", "port", 
    server.Addr, "URL", 
    "http://localhost"+server.Addr,
  )
  if err := server.ListenAndServe(); err != nil {
    fmt.Println("Error starting server: ", err)
  }
}
