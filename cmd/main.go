package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gvidow/YourVoiceWeb/pkg/cloud"
	ws "github.com/gvidow/YourVoiceWeb/pkg/websocket"
)

func main() {
	// read config
	cfg, err := ReadAllConfig()
	if err != nil {
		log.Fatal(err)
	}
	servCfg, cloudCfg, GPTCfg := GetServerConfig(cfg), GetYandexCloudConfig(cfg), GetChatGPTConfig(cfg)

	//create auto update config for yandex cloud (for workin SpeecKit)
	cloudConfig := cloud.NewCloudConfig(cloudCfg.oAuthToken, cloudCfg.folderId)

	err = cloudConfig.SrartAutoUpdateCloudConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer cloudConfig.Stop()

	//server
	serv := ws.NewWebSocketServer("", cloudConfig, GPTCfg.token)

	servMux := http.NewServeMux()

	mainPage, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	servMux.HandleFunc("/main", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(mainPage)
		if err != nil {
			log.Println(err)
		}
	}))
	dir := http.Dir("static")
	fs := http.FileServer(dir)
	servMux.Handle("/static/", http.StripPrefix("/static/", fs))
	servMux.Handle("/ws", serv)
	servMux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/main", http.StatusFound) }))
	if err = http.ListenAndServe(servCfg.host+":"+servCfg.port, servMux); err != nil {
		log.Fatal(err)
	}
}
