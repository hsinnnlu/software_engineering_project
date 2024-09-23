package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Successfully connected to database!")
}

func GetUserById(db *sql.DB, user_id string) (*models.User, error) {

	user := models.User{}
	query := "SELECT user_id, password_hash, user_permission FROM users WHERE user_id = ?"
	err := DB.QueryRow(query, user_id).Scan(&user.Id, &user.Password_hash, &user.Permission)
	if err != nil {
		// 如果找不到使用者
		if err == sql.ErrNoRows {
			fmt.Println("user not found")
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil

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

// 將密碼使用 SHA256
func getHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
