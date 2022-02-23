package telegram

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/persistence/entities/user"
	"GoBudgetBot/telegram/handlers"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

const (
	unrecognizedCommand = "Unrecognized command, please try again"
)

func Start(t string) {
	if len(t) == 0 {
		panic(fmt.Sprintf("Parameter %s is required, but was not passed", constants.TelegramToken))
	}

	bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	listen(bot)

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

/*
* Listens to messages from the Telegram channel, and responds to them
 */
func listen(botapi *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botapi.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			log.Fatal("Received nil message")
			return
		}

		fmt.Printf("Received message : " + update.Message.Text)

		u, err := user.GetByUserIdAndChatId(strconv.FormatInt(update.Message.From.ID, 10), strconv.FormatInt(update.Message.Chat.ID, 10))

		if err != nil {
			// user does not exist
			user.Create(user.FromMessage(update.Message))
		}

		handler, err := handlers.InitHandler(update.Message, u)

		var responseMessage string

		if err != nil {
			responseMessage = unrecognizedCommand
		} else {
			handler.Process()
			responseMessage = handler.GetResponseMessage()
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)
		response.ReplyToMessageID = update.Message.MessageID
		_, _ = botapi.Send(response)
	}
}
