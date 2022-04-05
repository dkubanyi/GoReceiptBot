package context

import (
	"GoBudgetBot/internal/domain/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotContext struct {
	Message *tgbotapi.Message
	User    *user.User
}
