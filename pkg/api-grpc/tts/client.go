package tts

import (
	"context"
	"fmt"
	"log"

	"github.com/gvidow/YourVoiceWeb/pkg/api-grpc/cloud"
	yatts "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/tts/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const addr = "tts.api.cloud.yandex.net:443"

type TextToSpeechClietn struct {
	synthesizer yatts.SynthesizerClient
	sendChan    chan string
	recvChan    chan []byte
	clientChan  chan yatts.Synthesizer_UtteranceSynthesisClient
	cfg         *cloud.CloudConfig
}

func (ttsc *TextToSpeechClietn) Start() error {
	iamToken, err := ttsc.cfg.GetIAMToken()
	if err != nil {
		return err
	}
	folderId, err := ttsc.cfg.GetFolderId()
	if err != nil {
		return err
	}

	go func() {
		ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", "Bearer "+iamToken, "x-folder-id", folderId)
		for text := range ttsc.sendChan {
			c, err := ttsc.synthesizer.UtteranceSynthesis(ctx, &yatts.UtteranceSynthesisRequest{Utterance: &yatts.UtteranceSynthesisRequest_Text{Text: text}})
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(4)
			ttsc.clientChan <- c
			fmt.Println(45)
			log.Println("YES")
		}
		close(ttsc.clientChan)
	}()
	go func() {
		fmt.Println(6)
		for cc := range ttsc.clientChan {
			fmt.Println(67)
			log.Println("NO")
			res, err := cc.Recv()
			if err != nil {
				log.Println(err)
				break
			}
			ttsc.recvChan <- res.AudioChunk.Data
		}
		close(ttsc.recvChan)
	}()
	return nil
}

func (ttsc *TextToSpeechClietn) GetSendChan() chan<- string {
	return ttsc.sendChan
}

func (ttsc *TextToSpeechClietn) GetRecvChan() <-chan []byte {
	return ttsc.recvChan
}

func NewTextToSpeechClient(cfg *cloud.CloudConfig) (*TextToSpeechClietn, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}
	cl := yatts.NewSynthesizerClient(conn)
	return &TextToSpeechClietn{
		synthesizer: cl, recvChan: make(chan []byte),
		sendChan:   make(chan string, 100),
		clientChan: make(chan yatts.Synthesizer_UtteranceSynthesisClient),
		cfg:        cfg,
	}, nil
}