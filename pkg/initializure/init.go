package initializure

import (
	"database/sql"
)

func ConnectToDb() {
	DB, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
}
