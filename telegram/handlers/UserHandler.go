package handlers

import (
	"GoBudgetBot/models"
	"errors"
	"fmt"
	"log"
)

/**
* This handler is responsible for processing user updates
 */
type userHandler struct {
	context models.BotContext
}

func (h *userHandler) IsResponsible() bool {
	t := h.context.Message.Text
	return t == CommandMe || t == CommandDeleteMe
}

func (h *userHandler) Process() error {
	if h.context.Message.Text == CommandDeleteMe {
		// TODO transaction
		if err := models.DeleteReceiptsByUserId(h.context.User); err != nil {
			log.Printf("Failed to delete receipts for user: %v", err)
			return errors.New("failed to delete receipts for user")
		}

		if _, err := models.DeleteUserById(h.context.User.Id); err != nil {
			log.Printf("Failed to delete user: %v", err)
			return errors.New("failed to delete user")
		}
	}

	return nil
}

func (h *userHandler) GetResponseMessage() string {
	if h.context.Message.Text == CommandDeleteMe {
		return fmt.Sprintf("Success.")
	}

	return fmt.Sprintf(
		"You are logged in as %s. Your account was created at %s.",
		h.context.User.Username,
		h.context.User.CreatedAt.Format("02.01.2006 15:04"),
	)
}
