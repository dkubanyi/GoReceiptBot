package main

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/telegram"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	telegram.Start(os.Getenv(constants.TELEGRAM_TOKEN))
}
