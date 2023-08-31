package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud/api-grpc/asr"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (s *Service) Asr(c *gin.Context) {
	s.log.Info("request for asr")
	defer s.log.Info("asr stop")
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")
	go func() {
		for i := 0; i < 10; i++ {
			for j := 0; j < 20; j++ {
				err := wsjson.Write(context.Background(), conn, &asr.Response{
					Text: "one ",
				})
				if err != nil {
					return
				}
				time.Sleep(20 * time.Millisecond)
			}
			err := wsjson.Write(context.Background(), conn, &asr.Response{
				Text: "two ",
				Fix:  true,
			})
			if err != nil {
				return
			}
		}
		wsjson.Write(context.Background(), conn, &asr.Response{
			Text:   "3 ",
			Fix:    true,
			Finish: true,
		})
	}()
	k := 0
	for {
		mt, m, err := conn.Read(context.Background())
		if err != nil {
			break
		}
		switch mt {
		case websocket.MessageBinary:
			k++
		case websocket.MessageText:
			fmt.Println(string(m), len(string(m)))
			fmt.Println("Pack: ", k)
			return
		}

	}
}
