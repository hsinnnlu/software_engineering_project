package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

var UserData map[string]string

type User struct {
	user_id  string
	password string
}

var DB *sql.DB

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
}

func init() {
	// UserData = map[string]string{
	// 	"test": "test",
	// }
	type User struct {
		user_id  string
		password string
	}
}

// CheckUserExists checks if a user exists in the database
func CheckUserExists(user_id string) (bool, error) {
	var user User
	err := DB.QueryRow("SELECT user_id, password FROM Users WHERE user_id = ?", user_id).Scan(&user.user_id, &user.password)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Auth authenticates a user by username and password
func Auth(username string, password string) error {
	var user User
	err := DB.QueryRow("SELECT user_id, password FROM Users WHERE user_id = ?", username).Scan(&user.user_id, &user.password)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user does not exist")
		}
		return err
	}

	if user.password != password {
		return errors.New("password is incorrect")
	}

	return nil
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginAuth(c *gin.Context) {
	var (
		user_id  string
		password string
	)
	if in, isExist := c.GetPostForm("user_id"); isExist && in != "" {
		user_id = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入使用者名稱"),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入密碼名稱"),
		})
		return
	}
	if err := Auth(user_id, password); err == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"success": "登入成功",
		})
		return
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err,
		})
		return
	}
}

func main() {

	InitDB("./database.db")
	defer DB.Close()

	fmt.Println("Data inserted successfully!")

	server := gin.Default()
	server.LoadHTMLGlob("./login.html")
	//設定靜態資源的讀取)
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.Run(":8888")
}
