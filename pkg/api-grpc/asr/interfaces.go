package asr

type AutomaticSpeechRecognizer interface {
	GetSendChan() chan<- []byte
	GetRecvChan() <-chan string
}

// type AudioChunker interface {
// 	ReadAll() []byte
// 	io.Writer
// 	IsLastChunk() bool
// }

// type RecognitionElement interface {
// 	fmt.Stringer
// 	IsFinalVersion() bool
// }
