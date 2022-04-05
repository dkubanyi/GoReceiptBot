package main

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/models"
	"GoBudgetBot/telegram"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// test connection
	models.DB = createConnection()

	telegram.Start(os.Getenv(constants.TelegramToken))
}

func createConnection() *sql.DB {
	// Open the connection
	log.Println("Attempting to establish connection...")
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	if err = db.Ping(); err != nil {
		log.Println("Connection failed")
		panic(err)
	}

	log.Println("Successfully connected!")

	// return the connection
	return db
}
