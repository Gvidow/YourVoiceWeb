package websocket

import (
	"context"
	"net"
	"net/http"

	"github.com/gvidow/YourVoiceWeb/logger"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud"
	ws "nhooyr.io/websocket"
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

func (wss *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wss.serv.Handler.ServeHTTP(w, r)
}

// /////////////////////////////////////////////////////=NEW=/////////////////////////////////////////
type textAnswerGenerator interface {
	SendQuestion(ctx context.Context, question string) (<-chan string, error)
}

type speechTextConverter interface {
	StartAutomaticSpeechRecognition(ctx context.Context)
	StartSpeechSynthesis(ctx context.Context)
}

type WebSocketService struct {
	log    *logger.Logger
	ansGen textAnswerGenerator
}

func NewWebSocketService(log *logger.Logger) *WebSocketService {
	return &WebSocketService{log: log}
}

func (s *WebSocketService) VoiceQuestionStrimingRecognition(ctx context.Context, conn *ws.Conn) {
	// s.clo
	go func() {
		for {
			ms, b, err := conn.Read(ctx)
			if err != nil {
				s.log.Error("reading voice question from ws connect: " + err.Error())
				break
			}
			if ms == ws.MessageText && string(b) == "DONE" {
				break
			}
			// sendChan <- b
		}
		// close(sendChan)
	}()

	// for res := range recvChan {
	// 	if res.Fix {
	// 		question.WriteString(res.Text)
	// 	}
	// 	err := wsjson.Write(ctx, conn, res)
	// 	if err != nil {
	// 		s.log.Error("write response text question in ws connect: " + err.Error())
	// 	}
	// }
}

func (s *WebSocketService) ReadQuestionAndChatID(ctx context.Context, conn *ws.Conn) {

}
