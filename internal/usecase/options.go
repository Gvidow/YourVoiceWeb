package usecase

type Option interface {
	apply(*Usecase)
}

type optionFunc func(*Usecase)

func (f optionFunc) apply(u *Usecase) {
	f(u)
}

func SetChatRepository(c ChatRepository) Option {
	return optionFunc(func(u *Usecase) {
		u.chats = c
	})
}
