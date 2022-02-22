package handlers

import (
	"GoBudgetBot/persistence/entities/user"
	"fmt"
)

/**
* This handler is responsible for processing telegramUser updates
 */
type userHandler struct {
	text         string
	telegramUser user.User
}

func (h *userHandler) IsResponsible() bool {
	return h.text == "/me"
}

func (h *userHandler) Process() {
}

func (h *userHandler) GetResponseMessage() string {
	return fmt.Sprintf("You are logged in as %s. Your account was created at %s.", h.telegramUser.Username, h.telegramUser.CreatedAt.Format("02.01.2006"))
}
