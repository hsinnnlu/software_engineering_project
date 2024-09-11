package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"

	"github.com/hsinnnlu/software_engineering_project/models"
)

var DB *sql.DB // 全局變量, 小寫db是外地注入變量

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 确保连接是可用的
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Successfully connected to database!")
}

func GetUserById(db *sql.DB, user_id string) (*models.User, error) {

	user := models.User{}
	query := "SELECT user_id, password_hash FROM users WHERE user_id = ?"
	err := DB.QueryRow(query, user_id).Scan(&user.Id, &user.Password_hash)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func VerifyPassword(db *sql.DB, id string, password string) (bool, error) {
	var user, err = GetUserById(DB, id)
	if err != nil {
		return false, err
	}
	return user.Password_hash == password, nil
}

func TestDB(c *gin.Context, db *sql.DB, user_id string, password string) error {
	// 准备插入数据
	stmt, err := db.Prepare("INSERT INTO Users(user_id, password_hash) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 插入数据
	_, err = stmt.Exec(user_id, password)
	if err != nil {
		return err
	}

	fmt.Print("Insert user success\n")

	if err != nil {
		c.JSON(200, gin.H{
			"message": "error",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "success",
		})
	}
	return nil
}
