package handlers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ResponseHandler interface {
	IsResponsible() bool
	Process()
	GetResponseMessage() string
}

func InitHandler(message *tgbotapi.Message) (ResponseHandler, error) {
	handlers := []ResponseHandler{
		startHandler{text: message.Text},
		imageHandler{text: message.Text, image: message.Photo},
	}

	for _, h := range handlers {
		if h.IsResponsible() {
			return h, nil
		}
	}

	return nil, errors.New("handler not implemented")
}
