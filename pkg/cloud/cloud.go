package cloud

type TokenGetter interface {
	GetFolderId() (string, error)
	GetIAMToken() (string, error)
}

type AutomaticSpeechRecognizer interface{}

type cloudService struct {
	cfg TokenGetter
	asr AutomaticSpeechRecognizer
}
