package handlers

import "GoBudgetBot/models"

/**
* This handler is responsible for handling the initial "/start" command
 */
type startHandler struct {
	context models.BotContext
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
