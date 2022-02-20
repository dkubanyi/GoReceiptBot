package entities

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func CreateConnection() *sql.DB {
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
