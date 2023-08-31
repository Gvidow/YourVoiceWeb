package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gvidow/YourVoiceWeb/internal/pkg/repository/chat"
	"github.com/gvidow/YourVoiceWeb/internal/pkg/service"
	"github.com/gvidow/YourVoiceWeb/internal/routing"
	"github.com/gvidow/YourVoiceWeb/internal/usecase"
	"github.com/gvidow/YourVoiceWeb/logger"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud"
)

func main() {
	log, err := logger.New(
		logger.FormatTime(time.RFC3339),
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer log.Sync()

	tmpl := template.Must(template.ParseGlob(templates))

	cfg, err := readConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, err := NewConnectToMongo(ctx, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		err = db.Client().Disconnect(context.Background())
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	u := usecase.New(
		usecase.SetChatRepository(chat.New(db)),
	)

	cloudConfig := cloud.NewCloudConfig(cfg)
	_ = cloudConfig

	// err = cloudConfig.SrartAutoUpdateCloudConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cloudConfig.Stop()

	s := service.New(log, u, tmpl)

	r := routing.New(cfg, log)
	r.ProduceRouting(s)
	if err = r.Run(); err != nil {
		log.Fatal(err.Error())
	}

	// read config
	// cfg, err := ReadAllConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// servCfg, cloudCfg, GPTCfg := GetServerConfig(cfg), GetYandexCloudConfig(cfg), GetChatGPTConfig(cfg)

	// //create auto update config for yandex cloud (for workin SpeecKit)
	// cloudConfig := cloud.NewCloudConfig(cloudCfg.oAuthToken, cloudCfg.folderId)

	// err = cloudConfig.SrartAutoUpdateCloudConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cloudConfig.Stop()

	//server
	// serv := ws.NewWebSocketServer("", cloudConfig, GPTCfg.token)

	// servMux := http.NewServeMux()

	// mainPage, err := os.ReadFile("index.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// servMux.HandleFunc("/main", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	_, err := w.Write(mainPage)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }))
	// dir := http.Dir("static")
	// fs := http.FileServer(dir)
	// servMux.Handle("/static/", http.StripPrefix("/static/", fs))
	// servMux.Handle("/ws", serv)
	// servMux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/main", http.StatusFound) }))

	// _ = serv

	//GIN
	// gin.SetMode(gin.DebugMode)
	// r := gin.New()
	// r.Static("/static", "./static")

	// tmp, err := template.ParseGlob("templates/*")
	// // fmt.Println(tmp.Tree)

	// // tmp, err = tmp.New("fdf").Parse("<h1> {{.Text}} <h1>")
	// // fmt.Println(tmp.Tree)
	// fmt.Println(err)
	// r.GET("/main", func(ctx *gin.Context) {
	// 	ctx.Render(http.StatusOK, render.HTML{
	// 		Template: tmp,
	// 		Name:     "index.html",
	// 		Data:     map[string]string{"Text": "<h1>question from aaaa</h1>"},
	// 	})
	// })

	// if err := r.Run(":8080"); err != nil {
	// 	log.Fatal(err.Error())
	// }
}
