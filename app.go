package main

import (
	"fmt"
	"net/http"
	"simple-chat/socket"

	"golang.org/x/net/websocket"
	Server "simple-chat/server"
)

func main() {
	go func() {
		server := socket.NewServer()
		client := Server.New()
		http.Handle("/ws", websocket.Handler(server.HandleWS))
		http.Handle("/", client)
		fmt.Println("Server started at :11000")
		fmt.Println("Press enter to exit")
		err := http.ListenAndServe(":11000", nil)
		if err != nil {
			fmt.Println("Error starting server: ", err)
		}
	}()
	fmt.Scanln()
}
