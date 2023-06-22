package websocket

import (
	"net"
	"net/http"

	"github.com/gvidow/YourVoiceWeb/pkg/api-grpc/cloud"
)

type WebSocketServer struct {
	Addr string
	l    net.Listener
	serv *http.Server
}

func NewWebSocketServer(addr string, cfg *cloud.CloudConfig, tokenGPT string) *WebSocketServer {
	return &WebSocketServer{
		Addr: addr,
		serv: &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { mainHandle(w, r, cfg, tokenGPT) })},
	}
}

func (wss *WebSocketServer) Run() error {
	var err error
	wss.l, err = net.Listen("tcp", wss.Addr)
	if err != nil {
		return err
	}
	return wss.serv.Serve(wss.l)
}
