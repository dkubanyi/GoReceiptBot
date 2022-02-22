package handlers

import (
	"GoBudgetBot/persistence/entities/user"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	defaultMessage = "Welcome to BudgetBot. Try one of the following commands:\n" +
		"/start --> display this message\n" +
		"/me --> show your information"
)

type ResponseHandler interface {
	IsResponsible() bool
	Process()
	GetResponseMessage() string
}

func InitHandler(message *tgbotapi.Message, u user.User) (ResponseHandler, error) {
	handlers := []ResponseHandler{
		&startHandler{text: message.Text},
		&imageHandler{text: message.Text, image: message.Photo},
		&userHandler{
			text:         message.Text,
			telegramUser: u,
		},
	}

	for _, h := range handlers {
		if h.IsResponsible() {
			return h, nil
		}
	}

	return nil, errors.New("handler not implemented")
}
