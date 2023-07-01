package websocket

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gvidow/YourVoiceWeb/pkg/cloud/api-grpc/asr"
	"github.com/gvidow/YourVoiceWeb/pkg/gpt"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func VoiceQuestionStrimingRecognition(ctx context.Context, conn *ws.Conn,
	question *strings.Builder, sendChan chan<- []byte, recvChan <-chan *asr.Response) {

	go func() {
		for {
			ms, b, err := conn.Read(ctx)
			if err != nil {
				log.Println(fmt.Errorf("mainHandle: read connection: %w", err))
				break
			}
			if ms == ws.MessageText && string(b) == "DONE" {
				log.Println("Запись получена")
				break
			}
			sendChan <- b
		}
		close(sendChan)
	}()

	for res := range recvChan {
		if res.Fix {
			question.WriteString(res.Text)
		}
		err := wsjson.Write(ctx, conn, res)
		if err != nil {
			log.Println(fmt.Errorf("VoiceQuestionStrimingRecognition: write response: %w", err))
		}
	}
	log.Printf("Распознан вопрос: %s\n", question.String())
}

func ResponseFromGPTAndSendToTTS(ctx context.Context, conn *ws.Conn, gpt *gpt.ChatGPT,
	question string, sendChan chan<- string) {

	chGPT := gpt.Ask(question)
	var answer = new(strings.Builder)
	var delta = new(strings.Builder)
	for u := range chGPT {
		answer.WriteString(u)
		delta.WriteString(u)
		conn.Write(ctx, ws.MessageText, []byte(u))
		r, _ := utf8.DecodeRuneInString(u)
		if delta.Len() > 10 && unicode.IsPunct(r) {
			sendChan <- delta.String()
			delta.Reset()
		}

	}
	if len(delta.String()) > 0 {
		sendChan <- delta.String()
	}
	close(sendChan)
	log.Printf("Получен ответ от ChatGPT: %s\n", answer.String())
}

func SendChunkedRecord(ctx context.Context, conn *ws.Conn, recvChan <-chan []byte) {
	internalChan := make(chan struct{})
	go func() {
		for _, ok := <-internalChan; ok; _, ok = <-internalChan {
			_, _, err := conn.Read(ctx)
			if err != nil {
				log.Printf("function syncing: %v\n", err)
			}
		}
	}()

	for a := range recvChan {
		internalChan <- struct{}{}
		conn.Write(ctx, ws.MessageBinary, a)
	}
	close(internalChan)
}
