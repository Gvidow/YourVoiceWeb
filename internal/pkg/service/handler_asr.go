package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

var WaitingTimeReceiveQuestion time.Duration = 10 * time.Second

func (s *Service) Asr(c *gin.Context) {
	s.log.Info("request for asr")
	defer s.log.Info("asr stop")
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")
	s.wsService.VoiceQuestionStrimingRecognition(c.Request.Context(), conn)

	ctx, cancel := context.WithTimeout(c.Request.Context(), WaitingTimeReceiveQuestion)
	defer cancel()
	s.wsService.ReadQuestionAndChatID(ctx, conn)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		for j := 0; j < 20; j++ {
	// 			err := wsjson.Write(context.Background(), conn, &asr.Response{
	// 				Text: "one ",
	// 			})
	// 			if err != nil {
	// 				return
	// 			}
	// 			time.Sleep(20 * time.Millisecond)
	// 		}
	// 		err := wsjson.Write(context.Background(), conn, &asr.Response{
	// 			Text: "two ",
	// 			Fix:  true,
	// 		})
	// 		if err != nil {
	// 			return
	// 		}
	// 	}
	// 	wsjson.Write(context.Background(), conn, &asr.Response{
	// 		Text:   "3 ",
	// 		Fix:    true,
	// 		Finish: true,
	// 	})
	// }()
	// k := 0
	// for {
	// 	mt, m, err := conn.Read(context.Background())
	// 	if err != nil {
	// 		break
	// 	}
	// 	switch mt {
	// 	case websocket.MessageBinary:
	// 		k++
	// 	case websocket.MessageText:
	// 		fmt.Println(string(m), len(string(m)))
	// 		fmt.Println("Pack: ", k)
	// 		return
	// 	}

	// }
}
