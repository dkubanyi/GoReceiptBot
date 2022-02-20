package user

import "GoBudgetBot/persistence/entities"

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func GetById(id int64) {
	db := entities.CreateConnection()

	defer db.Close()
}
