package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ResponseHandler interface {
	isResponsible() bool
	process()
	getResponseMessage() string
}

func InitHandler(message *tgbotapi.Message) (ResponseHandler, error) {
	handlers := []ResponseHandler{
		startHandler{text: message.Text},
		imageHandler{text: message.Text, image: message.Photo},
	}

	for _, h := range handlers {
		if h.isResponsible() {
			return h, nil
		}
	}

	return nil, errors.New("handler not implemented")
}

type startHandler struct {
	text string
}

func (h startHandler) isResponsible() bool {
	return h.text == "/start"
}

func (h startHandler) process() {
	// TODO save into DB, etc
}

func (h startHandler) getResponseMessage() string {
	return defaultMessage
}

type imageHandler struct {
	text  string
	image []tgbotapi.PhotoSize
}

func (h imageHandler) isResponsible() bool {
	return len(h.image) != 0
}

func (h imageHandler) process() {
}

func (h imageHandler) getResponseMessage() string {
	return "bye bye!"
}
