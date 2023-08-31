package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	stt "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v3"
)

var addr string = ":9092"

type FakeServer struct {
	stt.UnimplementedRecognizerServer
}

const filename = "fakes/fakeData/text/default.txt"

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	stt.RegisterRecognizerServer(s, &FakeServer{})
	log.Println("StartFakeASRServer")
	if err = s.Serve(l); err != nil {
		log.Fatal(err)
	}
}

func (fs *FakeServer) RecognizeStreaming(r stt.Recognizer_RecognizeStreamingServer) error {
	log.Println("Call FakeServer.RecognizeStreaming")
	countRecvPack := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := r.Recv()
			if err != nil {
				log.Println(fmt.Errorf("fakeYaCloudServer: recognizeStreaming recv: %w", err))
				break
			}
			if _, ok := msg.GetEvent().(*stt.StreamingRequest_Eou); ok {
				break
			}
			countRecvPack++
		}
		log.Printf("end reading: %d packages\n", countRecvPack)
	}()
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		a := &stt.Alternative{}
		a.SetText("Произошла ошибка")
		f := &stt.FinalRefinement{}
		f.SetNormalizedText(&stt.AlternativeUpdate{Alternatives: []*stt.Alternative{a}})
		r.Send(&stt.StreamingResponse{Event: &stt.StreamingResponse_FinalRefinement{FinalRefinement: f}})
	} else {
		pkg := 1
		response := strings.Builder{}
		br := bufio.NewReader(file)
		for b, err := br.ReadByte(); err == nil; b, err = br.ReadByte() {
			if b == '\n' {
				err = response.WriteByte(' ')
			} else {
				err = response.WriteByte(b)
			}
			if err != nil {
				log.Println(err)
				break
			}
			if b == ' ' {
				err = r.Send(&stt.StreamingResponse{Event: makeResponsePartial(response.String())})
			} else if b == '\n' {
				err = r.Send(&stt.StreamingResponse{Event: makeResponseFinalRefinement(response.String())})
				response.Reset()
			} else {
				continue
			}
			if err != nil {
				log.Println(err)
			}
			log.Printf("Send PKG: %d\n", pkg)
			pkg++
			time.Sleep(time.Millisecond * 200)
		}
		err = r.Send(&stt.StreamingResponse{Event: makeResponseFinalRefinement(response.String())})
		if err != nil {
			log.Println(err)
		}
		log.Printf("Send PKG: %d\n", pkg)
	}
	wg.Wait()
	sc := &stt.StatusCode{}
	sc.SetCodeType(stt.CodeType_CLOSED)
	err = r.Send(&stt.StreamingResponse{Event: &stt.StreamingResponse_StatusCode{StatusCode: sc}})
	if err != nil {
		log.Println(err)
	}
	log.Println("Exit in ASR")
	return nil
}

func buildAlternativesWithText(text string) []*stt.Alternative {
	a := &stt.Alternative{}
	a.SetText(text)
	return []*stt.Alternative{a}
}

func makeResponsePartial(text string) *stt.StreamingResponse_Partial {
	return &stt.StreamingResponse_Partial{
		Partial: &stt.AlternativeUpdate{
			Alternatives: buildAlternativesWithText(text),
		},
	}
}

func makeResponseFinalRefinement(text string) *stt.StreamingResponse_FinalRefinement {
	f := &stt.FinalRefinement{}
	f.SetNormalizedText(&stt.AlternativeUpdate{
		Alternatives: buildAlternativesWithText(text),
	})
	return &stt.StreamingResponse_FinalRefinement{
		FinalRefinement: f,
	}
}
