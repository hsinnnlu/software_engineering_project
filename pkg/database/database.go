package database

import (
	"database/sql"
	"errors"

	model "github.com/ryan2156/software_engineering_project/pkg/models"
)

func SelectUser(DB *sql.DB, user_id string) (*model.User, error) {
	var user model.User

	var row = DB.QueryRow("SELECT * FROM Users WHERE user_id = ?", user_id)
	err := row.Scan(&user.Id, &user.Password_hash, &user.Email, &user.Permission)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return &user, nil

}
