package handlers

import (
	"GoBudgetBot/models/entities"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CommandStart    = "/start"
	CommandMe       = "/me"
	CommandDeleteMe = "/deleteMe"
	defaultMessage  = "Welcome to BudgetBot.\n" +
		"Try uploading a photo of a QR code, and I will do my best to process it!💪\n" +
		"Alternatively, try one of the following commands:\n" +
		CommandStart + " --> display this message\n" +
		CommandMe + " --> show your user information\n" +
		CommandDeleteMe + " --> delete all your data in BudgetBot. ⚠️ Warning ⚠️ This is irreversible!!!"
)

type ResponseHandler interface {
	IsResponsible() bool
	Process()
	GetResponseMessage() string
}

func InitHandler(message *tgbotapi.Message, u entities.User) (ResponseHandler, error) {
	handlers := []ResponseHandler{
		&startHandler{text: message.Text},
		&imageHandler{text: message.Text, image: message.Photo, user: u},
		&userHandler{
			text: message.Text,
			user: u,
		},
	}

	for _, h := range handlers {
		if h.IsResponsible() {
			return h, nil
		}
	}

	return nil, errors.New("handler not implemented")
}
