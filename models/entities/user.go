package entities

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	UserId    string    `json:"userId"`
	ChatId    string    `json:"chatId"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	db := CreateConnection()
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
	db := CreateConnection()
	defer db.Close()
	var user User

	row := db.QueryRow(`SELECT * FROM users WHERE user_id = $1 AND chat_id = $2`, userId, chatId)
	err := row.Scan(&user.Id, &user.UserId, &user.ChatId, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)

	switch err {
	case sql.ErrNoRows:
		return user, errors.New("no rows were returned")
	case nil:
		return user, nil
	default:
		return user, errors.New(fmt.Sprintf("unable to scan the row. %v", err))
	}
}

func Create(user User) (User, error) {
	db := CreateConnection()
	defer db.Close()

	var id string
	sqlStatement := "INSERT INTO users (id, user_id, chat_id, username, first_name, last_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id"

	err := db.QueryRow(
		sqlStatement,
		uuid.New(),
		user.UserId,
		user.ChatId,
		user.Username,
		user.FirstName,
		user.LastName,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	log.Printf("Inserted new user with user_id %v", id)
	return GetByUserIdAndChatId(user.UserId, user.ChatId)
}

func DeleteById(uuid uuid.UUID) (bool, error) {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "DELETE FROM users WHERE id=$1"

	_, err := db.Exec(sqlStatement, uuid)

	if err != nil {
		return false, errors.New(fmt.Sprintf("Unable to execute query. %v", err))
	}

	return true, nil
}
