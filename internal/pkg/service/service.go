package service

import (
	"html/template"

	"github.com/gvidow/YourVoiceWeb/internal/pkg/websocket"
	"github.com/gvidow/YourVoiceWeb/internal/usecase"
	"github.com/gvidow/YourVoiceWeb/logger"
)

type Service struct {
	log       *logger.Logger
	usecase   *usecase.Usecase
	tmpl      *template.Template
	wsService *websocket.WebSocketService
}

func New(log *logger.Logger, u *usecase.Usecase, w *websocket.WebSocketService, tmpl *template.Template) *Service {
	return &Service{
		log:       log,
		tmpl:      tmpl,
		usecase:   u,
		wsService: w,
	}
}
