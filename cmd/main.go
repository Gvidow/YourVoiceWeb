package main

import (
	"log"
	"net/http"

	"github.com/gvidow/YourVoiceWeb/pkg/api-grpc/cloud"
	ws "github.com/gvidow/YourVoiceWeb/pkg/websocket"
	"go.uber.org/config"
)

func main() {
	//read config
	configFile := config.File("config.yaml")
	cfg, err := config.NewYAML(configFile)
	if err != nil {
		log.Fatal(err)
	}

	port := cfg.Get("server.port").String()
	host := cfg.Get("server.host").String()
	addr := host + ":" + port

	folderId := cfg.Get("yandexCloud.folderId").String()
	oAuthToken := cfg.Get("yandexCloud.OAuthToken").String()
	tokenGPT := cfg.Get("gpt.token").String()

	//create auto update config for yandex cloud (for workin SpeecKit)
	cloudConfig := cloud.NewCloudConfig(oAuthToken, folderId)

	err = cloudConfig.SrartAutoUpdateCloudConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer cloudConfig.Stop()

	//file server
	dir := http.Dir(".")
	h := http.FileServer(dir)
	s := &http.Server{Addr: "localhost:8081", Handler: h}
	go func() { log.Fatal(s.ListenAndServe()) }()

	//websocket server
	serv := ws.NewWebSocketServer(addr, cloudConfig, tokenGPT)

	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}

}
