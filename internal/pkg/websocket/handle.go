package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gvidow/YourVoiceWeb/pkg/cloud"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud/api-grpc/asr"
	"github.com/gvidow/YourVoiceWeb/pkg/cloud/api-grpc/tts"
	"github.com/gvidow/YourVoiceWeb/pkg/gpt"
	ws "nhooyr.io/websocket"
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

	// Запуск ASR от SpeechKit
	asrClient, err := asr.NewAutomaticSpeechRecognitionClient(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	asrClient.StartSpeechRecognition()

	// VoiceQuestionStrimingRecognition слушает websocket и паралльно отправляет по нему распознанную речь,
	// вызов блокируется до окончания распознования, т.к. для всё равно для отправки ChatGPT нам нужен полный вопрос.
	// Слушание прекращается после получение сообщения: "DONE".
	ctxReadVoice, cancelReadVoice := context.WithTimeout(context.Background(), MaxTimeKeepConnection)
	defer cancelReadVoice()
	var question strings.Builder
	VoiceQuestionStrimingRecognition(ctxReadVoice, conn, &question, asrClient.GetSendChan(), asrClient.GetRecvChan())

	// Запуск TTS от SpeechKit
	ttsClient, err := tts.NewTextToSpeechClient(cfg)
	if err != nil {
		log.Println(err)
	}
	ttsClient.Start()

	ctxWriteResponse, cancelWriteResponse := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelWriteResponse()

	clientGPT := gpt.NewChatGPT(tokenGPT)

	// Чтение ответа от ChatGPT по технологии SSE и отправка полученных промежуточных ответов на озвучку
	go ResponseFromGPTAndSendToTTS(ctxWriteResponse, conn, clientGPT, question.String(), ttsClient.GetSendChan())

	// Отправка озвучки ответа по кусочкам, для проигрывания на клиенте параллельно получению ответа от ChatGPT
	SendChunkedRecord(ctxWriteResponse, conn, ttsClient.GetRecvChan())

	time.Sleep(10 * time.Second)
	log.Println("Close connection")
}
