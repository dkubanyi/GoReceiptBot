package main

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/persistence/entities"
	"GoBudgetBot/telegram"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	// test connection
	db := entities.CreateConnection()
	defer db.Close()

	telegram.Start(os.Getenv(constants.TelegramToken))
}
