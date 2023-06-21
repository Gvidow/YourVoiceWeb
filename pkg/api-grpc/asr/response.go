package asr

type Response struct {
	Text   string `json:"text"`
	Fix    bool   `json:"fix"`
	Finish bool   `json:"finish"`
}
