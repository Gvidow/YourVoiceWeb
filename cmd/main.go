package main

import (
	"fmt"
	"log"

	ws "github.com/gvidow/YourVoiceWeb/pkg/websocket"
)

func main() {
	log.Println("Start")
	defer log.Println("Finish")

	const host = "localhost"
	const port = 8080

	addr := fmt.Sprintf("%s:%d", host, port)

	serv := ws.NewWebSocketServer(addr)

	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}

}
