package models

type User struct {
	Id            string
	Password_hash string
	Email         string
	Permission    string
}
