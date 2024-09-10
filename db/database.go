package db

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"

	"github.com/hsinnnlu/software_engineering_project/models"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() (*sql.DB, error) {
	DB, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Println("无法打开数据库:", err)
		return nil, err
	}

	// 使用 Ping 确保数据库连接成功
	if err := DB.Ping(); err != nil {
		log.Println("无法连接到数据库:", err)
		return nil, err
	}

	return DB, nil
}

func GetUserById(user_id string) (*models.User, error) {
	// var user models.User

	// var row = DB.QueryRow("SELECT * FROM Users WHERE user_id = ?", user_id)
	// err := row.Scan(&user.Id, &user.Password_hash, &user.Email, &user.Permission)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, errors.New("user does not exist")
	// 	}
	// 	return nil, err
	// }
	// return &user, nil

	user := &models.User{}
	query := "SELECT id, username, password FROM users WHERE id = ?"
	err := DB.QueryRow(query, user_id).Scan(&user.Id, &user.Password_hash, &user.Email, &user.Permission)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func VerifyPassword(Db *sql.DB, id string, password string) (bool, error) {
	var user, err = GetUserById(id)
	if err != nil {
		return false, err
	}
	return user.Password_hash == password, nil
}

func TestDB(c *gin.Context) {
	_, err := GetUserById("test")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "error",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "success",
		})
	}

}
