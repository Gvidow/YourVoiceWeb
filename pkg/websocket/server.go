package websocket

import (
	"net"
	"net/http"
)

type WebSocketServer struct {
	Addr string
	l    net.Listener
	serv *http.Server
}

func NewWebSocketServer(addr string) *WebSocketServer {
	return &WebSocketServer{
		Addr: addr,
		serv: &http.Server{Handler: http.HandlerFunc(mainHandle)},
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
