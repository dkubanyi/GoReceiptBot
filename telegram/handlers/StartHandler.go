package handlers

import "GoBudgetBot/telegram"

/**
* This handler is responsible for handling the initial "/start" command
 */
type startHandler struct {
	text string
}

func (h startHandler) IsResponsible() bool {
	return h.text == "/start"
}

func (h startHandler) Process() {
	// TODO save into DB, etc
}

func (h startHandler) GetResponseMessage() string {
	return telegram.DefaultMessage
}
