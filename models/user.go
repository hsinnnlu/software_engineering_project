package models

type User struct {
	Id            string
	Password_hash string
	Email         string
	Permission    string // 1: 學生, 2: 管理者, 3: 老師
}
