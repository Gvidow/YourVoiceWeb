package main

import (
	"io/ioutil"
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

	//server
	serv := ws.NewWebSocketServer("", cloudConfig, tokenGPT)

	servMux := http.NewServeMux()

	mainPage, err := ioutil.ReadFile("index.html")
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
	if err = http.ListenAndServe(addr, servMux); err != nil {
		log.Fatal(err)
	}
}
