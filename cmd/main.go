package main

import (
	"fmt"
	"log"

	"github.com/gvidow/YourVoiceWeb/pkg/api-grpc/cloud"
	ws "github.com/gvidow/YourVoiceWeb/pkg/websocket"
)

func main() {
	log.Println("Start")
	defer log.Println("Finish")

	const host = "localhost"
	const port = 8080

	addr := fmt.Sprintf("%s:%d", host, port)

	cloudConfig := cloud.NewCloudConfig("", "")

	err := cloudConfig.SrartAutoUpdateCloudConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer cloudConfig.Stop()

	var tokenGPT = ""

	serv := ws.NewWebSocketServer(addr, cloudConfig, tokenGPT)

	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}

}
