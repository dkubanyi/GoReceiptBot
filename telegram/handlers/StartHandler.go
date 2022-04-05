package handlers

import (
	"GoBudgetBot/internal/domain/context"
)

/**
* This handler is responsible for handling the initial "/start" command
 */
type startHandler struct {
	context context.BotContext
}

func (h *startHandler) IsResponsible() bool {
	return h.context.Message.Text == "/start"
}

func (h *startHandler) Process() error {
	return nil
}

func (h *startHandler) GetResponseMessage() string {
	return defaultMessage
}
