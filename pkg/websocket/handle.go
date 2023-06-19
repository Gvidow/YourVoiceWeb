package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	var record []byte
	for {
		ms, b, err := conn.Read(ctxRead)
		if err != nil {
			log.Println(fmt.Errorf("mainHandle: read connection: %w", err))
			break
		}
		if ms == ws.MessageText && string(b) == "DONE" {
			break
		}
		record = append(record, b...)

		log.Println(ms, len(b))
	}
	log.Printf("Размер полученной записи: %d", len(record))
	file, err := os.Create("record.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	n, err := file.Write(record)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)

	log.Println("Close connection")
}
