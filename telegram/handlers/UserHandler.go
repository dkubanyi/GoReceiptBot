package handlers

import (
	"GoBudgetBot/models/entities/user"
	"fmt"
)

/**
* This handler is responsible for processing user updates
 */
type userHandler struct {
	text string
	user user.User
}

func (h *userHandler) IsResponsible() bool {
	return h.text == CommandMe || h.text == CommandDeleteMe
}

func (h *userHandler) Process() {
	if h.text == CommandDeleteMe {
		user.DeleteById(h.user.Id)
	}
}

func (h *userHandler) GetResponseMessage() string {
	if h.text == CommandDeleteMe {
		return fmt.Sprintf("Success.")
	}

	return fmt.Sprintf("You are logged in as %s. Your account was created at %s.", h.user.Username, h.user.CreatedAt.Format("02.01.2006 15:04"))
}
