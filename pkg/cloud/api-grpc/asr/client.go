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
			// data := chunk.ReadAll()
			data := chunk
			// log.Println("goroutin send chunk: new chunk: size: ", len(data))
			err := asrc.Send(&stt.StreamingRequest{Event: NewChuck(data)})
			if err != nil {
				log.Println(fmt.Errorf("StartSpeechRecognition: goroutin send chunk: send chunk: %w", err))
			}
			// if chunk.IsLastChunk() {
			// 	err := asrc.Send(&stt.StreamingRequest{Event: NewEou()})
			// 	if err != nil {
			// 		log.Println(fmt.Errorf("StartSpeechRecognition: goroutin send chunk: send eou: %w", err))
			// 	}
			// }
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
			// fmt.Println("==========================")
			res, err := asrc.Recv()
			if err != nil {
				log.Println(fmt.Errorf("StartSpeechRecognition: goroutin recv text: recv: %w", err))
			}
			switch ev := res.Event.(type) {
			case *stt.StreamingResponse_StatusCode:
				log.Printf("Status code: %s\n", ev.StatusCode.CodeType)
				// asrc.recvChan <- "Status code"
			case *stt.StreamingResponse_FinalRefinement:
				log.Println("Final Refinement")
				if len(ev.FinalRefinement.GetNormalizedText().Alternatives) > 0 {
					asrc.recvChan <- &Response{Text: ev.FinalRefinement.GetNormalizedText().Alternatives[0].Text, Fix: true, Finish: false}
				}
				select {
				case <-asrc.stateRecordChan:
					break loop
				default:
					continue
				}
			case *stt.StreamingResponse_EouUpdate:
				log.Println("EouUpdate")
			case *stt.StreamingResponse_Partial:
				log.Printf("Partial: LEN %d\n", len(ev.Partial.Alternatives))
				if len(ev.Partial.Alternatives) > 0 {
					asrc.recvChan <- &Response{Text: ev.Partial.Alternatives[0].Text, Fix: false, Finish: false}
				}
				// fmt.Println(ev)
			case *stt.StreamingResponse_Final:
				fmt.Println("Final")
				//break loop
			}
			time.Sleep(100 * time.Millisecond)
		}
		// res, err := asrc.Recv()
		// if err != nil {
		// 	log.Println(fmt.Errorf("StartSpeechRecognition: goroutin recv text: recv: %w", err))
		// }
		// if ev, ok := res.Event.(*stt.StreamingResponse_FinalRefinement); ok {
		// 	if text, ok := ev.FinalRefinement.Type.(*stt.FinalRefinement_NormalizedText); ok {
		// 		asrc.recvChan <- text.NormalizedText.Alternatives[0].Text
		// 		fmt.Printf("LEN %d\n", len(text.NormalizedText.Alternatives))
		// 	}
		// 	// fmt.Println(ev)
		// }
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
