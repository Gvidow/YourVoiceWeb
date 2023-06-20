package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gvidow/YourVoiceWeb/pkg/api-grpc/asr"
	ws "nhooyr.io/websocket"
)

const MaxTimeKeepConnection time.Duration = 4*time.Minute + 50*time.Second

func mainHandle(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Origin")
	conn, err := ws.Accept(w, r, nil)
	if err != nil {
		log.Println(fmt.Errorf("mainHandle: Accept: %w", err))
		return
	}

	log.Println("New connection")
	ctxRead, cancelRead := context.WithTimeout(context.Background(), MaxTimeKeepConnection)
	defer cancelRead()

	asrClient, err := asr.NewAutomaticSpeechRecognitionClient()
	if err != nil {
		log.Println(err)
		return
	}
	asrClient.StartSpeechRecognition()

	ch := asrClient.GetSendChan()
	// var record []byte
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
			// record = append(record, b...)
			ch <- b
			// log.Println(ms, len(b))
		}
		// file, err := os.Create("record.wav")
		// if err != nil {
		// 	log.Println(err)
		// }
		// defer file.Close()
		// file.Write(record)
		// fmt.Println("lll", len(record))
	}()
	// log.Printf("Размер полученной записи: %d", len(record))
	ctx := context.Background()
	for text := range asrClient.GetRecvChan() {
		// log.Println("H: ", text)
		conn.Write(ctx, ws.MessageText, []byte(text))
	}
	log.Println("Close connection")
}
