package user

import (
	"GoBudgetBot/persistence/entities"
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
	"strconv"
)

type User struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	ChatId    string `json:"chatId"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FromMessage(update *tgbotapi.Message) User {
	return User{
		UserId:    strconv.FormatInt(update.From.ID, 10),
		ChatId:    strconv.FormatInt(update.Chat.ID, 10),
		Username:  update.From.UserName,
		FirstName: update.From.FirstName,
		LastName:  update.From.LastName,
	}
}

func Listing() ([]User, error) {
	db := entities.CreateConnection()
	defer db.Close()

	var users []User

	sqlStatement := `SELECT * FROM users`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.UserId, &user.ChatId, &user.Username, &user.FirstName, &user.LastName)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)
	}

	return users, err
}

func GetByUserIdAndChatId(userId string, chatId string) (User, error) {
	db := entities.CreateConnection()
	defer db.Close()
	var user User

	row := db.QueryRow(`SELECT * FROM users WHERE user_id = $1 AND chat_id = $2`, userId, chatId)
	err := row.Scan(&user.Id, &user.UserId, &user.ChatId, &user.Username, &user.FirstName, &user.LastName)

	switch err {
	case sql.ErrNoRows:
		return user, errors.New("no rows were returned")
	case nil:
		return user, nil
	default:
		return user, errors.New(fmt.Sprintf("unable to scan the row. %v", err))
	}
}

func Create(user User) (string, error) {
	db := entities.CreateConnection()
	defer db.Close()

	var id string
	sqlStatement := "INSERT INTO users (id, user_id, chat_id, username, first_name, last_name) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id"

	err := db.QueryRow(sqlStatement, uuid.New(), user.UserId, user.ChatId, user.Username, user.FirstName, user.LastName).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	log.Printf("Inserted new user with user_id %v", id)
	return id, nil
}
