package telegram

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/models/entities"
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

	// set to true if you want to see update information (messages received in Telegram) in console
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
			log.Println("Received nil message")
			continue
		}

		u, err := entities.GetByUserIdAndChatId(strconv.FormatInt(update.Message.From.ID, 10), strconv.FormatInt(update.Message.Chat.ID, 10))

		if err != nil {
			// user does not exists
			u, _ = entities.CreateUser(entities.FromMessage(update.Message))
		}

		handler, err := handlers.InitHandler(update.Message, &u)

		var responseMessage string

		if err != nil {
			responseMessage = unrecognizedCommand
		} else {
			handler.Process()
			responseMessage = handler.GetResponseMessage()
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)
		response.ParseMode = tgbotapi.ModeHTML
		response.ReplyToMessageID = update.Message.MessageID
		response.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(handlers.CommandStart),
				tgbotapi.NewKeyboardButton(handlers.CommandMe),
				tgbotapi.NewKeyboardButton(handlers.CommandShowReceipts),
			),
		)

		_, _ = botapi.Send(response)
	}
}
