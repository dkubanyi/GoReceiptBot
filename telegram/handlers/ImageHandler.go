package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

/**
* This handler is responsible for processing updates containing images
 */
type imageHandler struct {
	text  string
	image []tgbotapi.PhotoSize
}

func (h imageHandler) IsResponsible() bool {
	return len(h.image) != 0
}

func (h imageHandler) Process() {
}

func (h imageHandler) GetResponseMessage() string {
	return defaultMessage
}
