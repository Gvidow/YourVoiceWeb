package gpt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const url = "https://api.openai.com/v1/chat/completions"

type ChatGPT struct {
	token string
}

func NewChatGPT(token string) *ChatGPT {
	return &ChatGPT{token: token}
}

func (cgpt *ChatGPT) Ask(question string) <-chan string {
	ch := make(chan string)
	go func() {
		strBodyRequest := fmt.Sprintf(`{
			"model": "gpt-3.5-turbo",
			"messages": [{"role": "user", "content": "%s"}],
			"stream": true
		}`, question)
		body := strings.NewReader(strBodyRequest)
		req, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Authorization", "")
		req.Header.Add("Content-Type", "application/json")
		req.ContentLength = int64(len(strBodyRequest))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(5)
		defer res.Body.Close()
		sc := bufio.NewScanner(res.Body)
		for sc.Scan() {
			err = sc.Err()
			if err != nil {
				log.Println(err)
				break
			}
			if len(sc.Text()) < 6 {
				continue
			}
			event := sc.Text()[6:]
			var ev Event
			err = json.Unmarshal([]byte(event), &ev)
			if err != nil {
				log.Println(err)
				break
			}
			if ev.Choices[0].FinishReason == "stop" {
				break
			}
			ch <- string(ev.Choices[0].Delta.Content)
		}
		close(ch)
	}()
	return ch
}
