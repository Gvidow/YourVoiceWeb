package websocket

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gvidow/YourVoiceWeb/pkg/cloud"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud/api-grpc/asr"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud/api-grpc/tts"
	"github.com/gvidow/YourVoiceWeb/pkg/gpt"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const MaxTimeKeepConnection time.Duration = 4*time.Minute + 50*time.Second

func mainHandle(w http.ResponseWriter, r *http.Request, cfg *cloud.CloudConfig, tokenGPT string) {
	conn, err := ws.Accept(w, r, nil)
	if err != nil {
		log.Println(fmt.Errorf("mainHandle: Accept: %w", err))
		return
	}
	defer conn.Close(ws.StatusNormalClosure, "")
	log.Println("New connection")
	ctxRead, cancelRead := context.WithTimeout(context.Background(), MaxTimeKeepConnection)
	defer cancelRead()

	asrClient, err := asr.NewAutomaticSpeechRecognitionClient(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	asrClient.StartSpeechRecognition()

	ch := asrClient.GetSendChan()
	var record []byte
	var question string
	go func() {
		for {
			ms, b, err := conn.Read(ctxRead)
			if err != nil {
				log.Println(fmt.Errorf("mainHandle: read connection: %w", err))
				close(ch)
				break
			}
			if ms == ws.MessageText && string(b) == "DONE" {
				log.Println("Запись получена")
				close(ch)
				break
			}
			record = append(record, b...)
			ch <- b
		}
	}()
	ctx, fff := context.WithTimeout(context.Background(), time.Hour)
	defer fff()
	for res := range asrClient.GetRecvChan() {
		if res.Fix {
			question += res.Text
		}
		wsjson.Write(ctx, conn, res)
	}
	internalChan := make(chan struct{})
	go func() {
		ctx, c := context.WithTimeout(context.Background(), time.Hour)
		defer c()
		for _, ok := <-internalChan; ok; _, ok = <-internalChan {
			t, b, err := conn.Read(ctx)
			log.Println(t, b)
			if err != nil {
				log.Printf("dop: %v\n", err)
			}
		}
		log.Println("OTL")
	}()
	ttsClient, err := tts.NewTextToSpeechClient(cfg)
	if err != nil {
		log.Println(err)
	}
	go func() {
		clientGPT := gpt.NewChatGPT(tokenGPT)
		chGPT := clientGPT.Ask(question)

		ctx, ccc := context.WithTimeout(context.Background(), time.Hour)
		defer ccc()
		var answer = bytes.NewBufferString("")
		var delta = bytes.NewBufferString("")
		for u := range chGPT {
			answer.WriteString(u)
			fmt.Println(u)
			delta.WriteString(u)
			conn.Write(ctx, ws.MessageText, []byte(u))
			r, _ := utf8.DecodeRuneInString(u)
			if delta.Len() > 10 && unicode.IsPunct(r) {
				fmt.Println(delta.String())
				ttsClient.GetSendChan() <- delta.String()
				delta.ReadString(0)
			}

		}
		if len(delta.String()) > 0 {
			ttsClient.GetSendChan() <- delta.String()
		}
		close(ttsClient.GetSendChan())
		log.Println(answer.String())
		log.Println("exit chatgpt")
	}()
	log.Println("AAAAAAAAAAAAAAAAAAAAAa")
	ttsClient.Start()
	fmt.Println("21212121212")
	// ttsClient.GetSendChan() <- answer.String()
	// close(ttsClient.GetSendChan())
	fmt.Println("2")
	nctx, p := context.WithTimeout(context.Background(), time.Hour)
	defer p()
	for a := range ttsClient.GetRecvChan() {
		internalChan <- struct{}{}
		conn.Write(nctx, ws.MessageBinary, a)
	}
	close(internalChan)
	time.Sleep(30 * time.Second)
	log.Println("Close connection")
}
