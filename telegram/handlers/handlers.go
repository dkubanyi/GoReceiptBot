package handlers

import (
	"GoBudgetBot/models"
	"GoBudgetBot/models/entities"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CommandStart        = "/start"
	CommandMe           = "/me"
	CommandQr           = "/qr"
	CommandShowReceipts = "/showReceipts"
	CommandDeleteMe     = "/deleteMe"
	defaultMessage      = "Welcome to BudgetBot.\n" +
		"Try uploading a photo of a QR code, and I will do my best to process it!💪\n" +
		"Alternatively, try one of the following commands:\n" +
		CommandStart + " --> display this message\n" +
		CommandMe + " --> show your user information\n" +
		CommandQr + " [text] --> submits [text] as a QR code\n" +
		CommandShowReceipts + " --> show your saved receipts\n" +
		"\nYou can also choose to delete all your data associated with this bot." +
		"\nUse " + CommandDeleteMe + " to do that. \n⚠️ Warning ⚠️ This is irreversible!!!"
)

var response string

type ResponseHandler interface {
	IsResponsible() bool
	Process() error
	GetResponseMessage() string
}

func InitHandler(message *tgbotapi.Message, u *entities.User) (ResponseHandler, error) {
	context := models.BotContext{
		Message: message,
		User:    u,
	}

	handlers := []ResponseHandler{
		&startHandler{context},
		&imageHandler{context},
		&qrHandler{context},
		&receiptHandler{context},
		&userHandler{context},
	}

	for _, h := range handlers {
		if h.IsResponsible() {
			return h, nil
		}
	}

	return nil, errors.New("handler not implemented")
}
