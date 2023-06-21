package gpt

type Event struct {
	Choices []Choice
}

type Choice struct {
	Delta struct {
		Content string
	}
	FinishReason string `json:"finish_reason"`
}
