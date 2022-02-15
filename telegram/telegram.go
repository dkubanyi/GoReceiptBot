package telegram

import (
	"GoBudgetBot/classes"
	"GoBudgetBot/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

var (
	tg *tgbotapi.BotAPI
)

const (
	PARAM_REQUIRED = "Parameter %s is required, but was not passed"
)

func Start(t string) {
	if len(t) == 0 {
		panic(fmt.Sprintf(PARAM_REQUIRED, constants.TELEGRAM_TOKEN))
	}

	bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	startPolling(bot)

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

/*
* Performs polling of messages from the Telegram channel, and responds to them
 */
func startPolling(botapi *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botapi.GetUpdatesChan(u)

	b := classes.New(&classes.Handlers{
		Response: responseHandler,
	})

	_ = b

	for update := range updates {
		if update.Message == nil {
			log.Fatal("Received nil message")
			return
		}

		outMsg := fmt.Sprintf("Hello! This bot is in development. You sent this message: %s`", update.Message.Text)
		log.Printf("[%s] %s", update.Message.From.UserName, outMsg)

		response := tgbotapi.NewMessage(update.Message.Chat.ID, outMsg)
		response.ReplyToMessageID = update.Message.MessageID

		botapi.Send(response)
	}
}

func responseHandler(msg classes.OutgoingMessage) {
	id, err := strconv.ParseInt(msg.Target, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	m := tgbotapi.NewMessage(id, msg.Message)
	_ = m

	tg.Send(&m)
}
