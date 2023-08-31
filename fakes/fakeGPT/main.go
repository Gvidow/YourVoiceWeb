package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var addr string = ":9091"

var ErrFailCastToFlusher = errors.New("не удалось привести к Flucher")

func main() {
	http.HandleFunc("/gpt", generateAnswerSSE)
	log.Println("StartFakeGPTServer")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func generateAnswerSSE(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Authorization header: %s, Content-Type: %s",
		r.Method, r.Header.Get("Authorization"), r.Header.Get("Content-Type"))

	d := json.NewDecoder(r.Body)
	var msg = &struct{ Messages []struct{ Content string } }{}
	err := d.Decode(msg)
	var question string
	if err != nil {
		log.Println(err)
		question = "Не удалось вычитать вопрос из тела запроса"
	} else {
		question = msg.Messages[0].Content
	}

	w.Header().Set("Content-Type", "text/event-stream")

	answer := []string{"Пр", "иве", "т", ",", " п", "оль", "зова", "тель", "!", " Ваш",
		" в", "опр", "о", "с", ":", " ", question, "."}

	for _, d := range answer {
		_, err = writeAndFlush(w, buildDelta(d))
		if err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	_, err = writeAndFlush(w, []byte(`data: {"choices":[{"delta": {"content": ""}, "finish_reason": "stop"}]}`+"\n"))
	if err != nil {
		log.Println(err)
	}

	_, err = writeAndFlush(w, []byte("data: DONE\n"))
	if err != nil {
		log.Println(err)
	}
	log.Print("DONE\n\n\n")
}

func writeAndFlush(w io.Writer, b []byte) (int, error) {
	n, err := w.Write(b)
	if err != nil {
		return n, fmt.Errorf("writeAndFlush: error write: %w", err)
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
		return n, nil
	}
	return n, ErrFailCastToFlusher
}

func buildDelta(delta string) []byte {
	return []byte(fmt.Sprintf(`data: {"choices":[{"delta": {"content": "%s"}, "finish_reason": ""}]}`+"\n", delta))
}
