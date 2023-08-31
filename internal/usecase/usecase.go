package usecase

import (
	"context"
	"errors"

	"github.com/gvidow/YourVoiceWeb/internal/pkg/repository/chat"
)

type ChatRepository interface {
	SelectAllByOrder(ctx context.Context) ([]chat.ChatDoc, error)
	DeleteMany(ctx context.Context, ids []string) (int, error)
	SaveSettings(ctx context.Context, id string, settings *chat.Setting) error
	SwapPlaces(ctx context.Context, id1, id2 string) error
	AddNewChat(ctx context.Context, title string) (string, error)
	EditChat(ctx context.Context, id string, newTitle string) error
}

var (
	ErrNotSetChatRepository = errors.New("the ChatRepository value is not set")
)

type Usecase struct {
	chats ChatRepository
}

func New(options ...Option) *Usecase {
	res := &Usecase{}
	for _, option := range options {
		option.apply(res)
	}
	return res
}

func (u *Usecase) ChatRepository() (ChatRepository, error) {
	if u.chats == nil {
		return nil, ErrNotSetChatRepository
	}
	return u.chats, nil
}
