package main

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/models"
	"GoBudgetBot/telegram"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	// test connection
	models.DB = createConnection()

	telegram.Start(os.Getenv(constants.TelegramToken))
}

func createConnection() *sql.DB {
	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// return the connection
	return db
}
