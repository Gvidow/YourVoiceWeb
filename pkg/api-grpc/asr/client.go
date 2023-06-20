package asr

import (
	"context"
	"fmt"
	"log"
	"time"

	stt "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const addr = "stt.api.cloud.yandex.net:443"

type AutomaticSpeechRecognitionClient struct {
	stt.Recognizer_RecognizeStreamingClient
	sendChan        chan []byte
	recvChan        chan string
	stateRecordChan chan struct{}
}

func (asrc *AutomaticSpeechRecognitionClient) GetSendChan() chan<- []byte {
	return asrc.sendChan
}

func (asrc *AutomaticSpeechRecognitionClient) GetRecvChan() <-chan string {
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
				fmt.Printf("Status code: %s\n", ev.StatusCode.CodeType)
				asrc.recvChan <- "Status code"
			case *stt.StreamingResponse_FinalRefinement:
				fmt.Println("Final Refinement")
				s := ""
				if len(ev.FinalRefinement.GetNormalizedText().Alternatives) > 0 {
					s = ev.FinalRefinement.GetNormalizedText().String()
				}
				asrc.recvChan <- "FR" + s
				select {
				case <-asrc.stateRecordChan:
					break loop
				default:
					continue
				}
			case *stt.StreamingResponse_EouUpdate:
				fmt.Println("EouUpdate")
				asrc.recvChan <- "EouUpdate"
			case *stt.StreamingResponse_Partial:
				fmt.Printf("Partial: LEN %d\n", len(ev.Partial.Alternatives))
				if len(ev.Partial.Alternatives) > 0 {
					asrc.recvChan <- ev.Partial.Alternatives[0].Text
				}
				// fmt.Println(ev)
			case *stt.StreamingResponse_Final:
				fmt.Println("Final")
				asrc.recvChan <- "Final"
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
		close(asrc.recvChan)
	}()
	return nil
}

func NewAutomaticSpeechRecognitionClient() (*AutomaticSpeechRecognitionClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}

	rc := stt.NewRecognizerClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", "", "x-folder-id", "")
	res, err := rc.RecognizeStreaming(ctx)
	if err != nil {
		return nil, err
	}
	return &AutomaticSpeechRecognitionClient{Recognizer_RecognizeStreamingClient: res,
		sendChan:        make(chan []byte),
		recvChan:        make(chan string),
		stateRecordChan: make(chan struct{})}, nil

}
