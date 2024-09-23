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

// 新增講座 ， timestamp格式還沒確定
func InsertLecture(c *gin.Context, lecture models.Lecture) error {

	// 准备插入数据
	query := "INSERT INTO Lectures(lecture_id, lecture_name, lecture_speaker, lecture_timestamp, lecture_manager, lecture_location) values(?, ?, ?, ?, ?, ?)"
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	fmt.Print("pass A\n")

	// 插入数据
	_, err = stmt.Exec(lecture.Id, lecture.Name, lecture.Speaker, lecture.Timestamp, lecture.Manager, lecture.Location)
	if err != nil {
		fmt.Printf("Insert lecture failed, err: %v\n", err)
		return err
	}

	fmt.Print("Insert lecture success\n")

	c.JSON(200, gin.H{
		"message": "success for insert lecture",
	})

	return nil
}
