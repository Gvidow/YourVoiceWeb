package asr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gvidow/YourVoiceWeb/pkg/cloud"
	stt "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const addr = "stt.api.cloud.yandex.net:443"

type AutomaticSpeechRecognitionClient struct {
	stt.Recognizer_RecognizeStreamingClient
	sendChan        chan []byte
	recvChan        chan *Response
	stateRecordChan chan struct{}
}

func (asrc *AutomaticSpeechRecognitionClient) GetSendChan() chan<- []byte {
	return asrc.sendChan
}

func (asrc *AutomaticSpeechRecognitionClient) GetRecvChan() <-chan *Response {
	return asrc.recvChan
}

func (asrc *AutomaticSpeechRecognitionClient) StartSpeechRecognition() error {
	err := asrc.Send(&stt.StreamingRequest{Event: NewSessionOptions()})
	if err != nil {
		return err
	}
	go func() {
		for chunk := range asrc.sendChan {
			err := asrc.Send(&stt.StreamingRequest{Event: NewChuck(chunk)})
			if err != nil {
				log.Println(fmt.Errorf("StartSpeechRecognition: goroutin send chunk: send chunk: %w", err))
			}
		}
		err := asrc.Send(&stt.StreamingRequest{Event: NewEou()})
		if err != nil {
			log.Println(fmt.Errorf("StartSpeechRecognition: goroutin send chunk: send eou: %w", err))
		}
		close(asrc.stateRecordChan)
	}()
	go func() {
		time.Sleep(time.Millisecond * 100)
	loop:
		for {
			res, err := asrc.Recv()
			if err != nil {
				log.Println(fmt.Errorf("StartSpeechRecognition: goroutin recv text: recv: %w", err))
			}
			switch ev := res.Event.(type) {
			case *stt.StreamingResponse_FinalRefinement:
				if len(ev.FinalRefinement.GetNormalizedText().Alternatives) > 0 {
					asrc.recvChan <- &Response{Text: ev.FinalRefinement.GetNormalizedText().Alternatives[0].Text, Fix: true, Finish: false}
				}
				select {
				case <-asrc.stateRecordChan:
					break loop
				default:
					continue
				}
			case *stt.StreamingResponse_Partial:
				if len(ev.Partial.Alternatives) > 0 {
					asrc.recvChan <- &Response{Text: ev.Partial.Alternatives[0].Text, Fix: false, Finish: false}
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
		asrc.recvChan <- &Response{Finish: true}
		close(asrc.recvChan)
	}()
	return nil
}

func NewAutomaticSpeechRecognitionClient(cfg *cloud.CloudConfig) (*AutomaticSpeechRecognitionClient, error) {
	iamToken, err := cfg.GetIAMToken()
	if err != nil {
		return nil, err
	}
	folderId, err := cfg.GetFolderId()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}
	rc := stt.NewRecognizerClient(conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", "Bearer "+iamToken, "x-folder-id", folderId)
	res, err := rc.RecognizeStreaming(ctx)
	if err != nil {
		return nil, err
	}
	return &AutomaticSpeechRecognitionClient{Recognizer_RecognizeStreamingClient: res,
		sendChan:        make(chan []byte),
		recvChan:        make(chan *Response),
		stateRecordChan: make(chan struct{})}, nil

}
