package models

import (
	"GoBudgetBot/models/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotContext struct {
	Message *tgbotapi.Message
	User    *entities.User
}
