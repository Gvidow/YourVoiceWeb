package main

import (
	"io"
	"log"
	"net"
	"os"

	tts "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/tts/v3"
	"google.golang.org/grpc"
)

var addr string = ":9093"

const filename = "fakes/fakeData/speech/speech.wav"

type FakeServer struct {
	tts.UnimplementedSynthesizerServer
}

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	tts.RegisterSynthesizerServer(s, &FakeServer{})
	log.Println("StartFakeTTSServer")
	if err = s.Serve(l); err != nil {
		log.Fatal(err)
	}
}

func (fs *FakeServer) UtteranceSynthesis(r *tts.UtteranceSynthesisRequest, s tts.Synthesizer_UtteranceSynthesisServer) error {
	log.Println("Call FakeServer.UtteranceSynthesis")
	text := r.GetText()
	log.Println("Текст на синтез: ", text)
	res := &tts.UtteranceSynthesisResponse{}
	audio := &tts.AudioChunk{}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		audio.SetData([]byte{})
	} else {
		audio.SetData(b)
	}
	res.SetAudioChunk(audio)
	s.Send(res)
	return nil
}
