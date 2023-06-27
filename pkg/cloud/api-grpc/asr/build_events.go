package asr

import "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v3"

func NewSessionOptions() *stt.StreamingRequest_SessionOptions {
	format := &stt.AudioFormatOptions{
		AudioFormat: &stt.AudioFormatOptions_RawAudio{
			RawAudio: &stt.RawAudio{
				SampleRateHertz: 48000, AudioChannelCount: 1,
			},
		},
	}
	return &stt.StreamingRequest_SessionOptions{
		SessionOptions: &stt.StreamingOptions{
			RecognitionModel: &stt.RecognitionModelOptions{
				AudioFormat: format,
			},
		},
	}
}

func NewChuck(data []byte) *stt.StreamingRequest_Chunk {
	return &stt.StreamingRequest_Chunk{Chunk: &stt.AudioChunk{Data: data}}
}

func NewEou() *stt.StreamingRequest_Eou {
	return &stt.StreamingRequest_Eou{Eou: &stt.Eou{}}
}
