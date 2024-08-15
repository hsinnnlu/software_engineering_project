package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"database/sql"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"

	_ "github.com/mattn/go-sqlite3"
)

type VerificationCode struct {
	Code       string
	Expiration time.Time
}

var UserData map[string]string
var loginAttempts map[string]int
var lockoutTime map[string]time.Time
var verifyCodes map[string]VerificationCode
var mu sync.Mutex

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
	verifyCodes = make(map[string]VerificationCode)
	loginAttempts = make(map[string]int)
	lockoutTime = make(map[string]time.Time)
}

func generateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

func sendMail(to, subject, body, from, password string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("mail.cs.thu.edu.tw", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func resendCodeHandler(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
		return
	}

	mu.Lock()
	verifyCode, exists := verifyCodes[requestBody.Email]
	mu.Unlock()

	if !exists || time.Now().After(verifyCode.Expiration) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "無此用戶的驗證碼"})
		return
	}

	// 發送郵件
	from := "nmg@cs.thu.edu.tw"
	password := "e04su3su;6"
	subject := "重送驗證碼"
	body := fmt.Sprintf("您的驗證碼是: %s", verifyCode.Code)

	err := sendMail(requestBody.Email, subject, body, from, password)
	if err != nil {
		log.Println("發送郵件失敗:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "發送郵件失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "驗證碼已重新發送"})
}

func sendCodeHandler(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
		return
	}

	// 生成驗證碼
	verifyCode := generateVerifyCode()

	// 保存驗證碼及其過期時間
	mu.Lock()
	verifyCodes[requestBody.Email] = VerificationCode{
		Code:       verifyCode,
		Expiration: time.Now().Add(2 * time.Minute), // 設置驗證碼過期時間
	}
	mu.Unlock()

	// 發送郵件
	from := "nmg@cs.thu.edu.tw"
	password := "e04su3su;6"
	subject := "驗證碼"
	body := fmt.Sprintf("您的驗證碼是: %s", verifyCode)

	err := sendMail(requestBody.Email, subject, body, from, password)
	if err != nil {
		log.Println("發送郵件失敗:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "發送郵件失敗"})
		return
	}

	// 設置 90 秒後檢查是否需要重送驗證碼
	time.AfterFunc(90*time.Second, func() {
		mu.Lock()
		defer mu.Unlock()

		// 如果驗證碼過期或用戶仍需驗證碼，則重新發送
		if code, exists := verifyCodes[requestBody.Email]; exists && time.Now().After(code.Expiration) {
			newCode := generateVerifyCode()
			verifyCodes[requestBody.Email] = VerificationCode{
				Code:       newCode,
				Expiration: time.Now().Add(5 * time.Minute),
			}
			sendMail(requestBody.Email, subject, fmt.Sprintf("您的新驗證碼是: %s", newCode), from, password)
		}
	})

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "驗證碼已發送"})
}

func verifyCodeHandler(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	println(requestBody.Email)
	println(requestBody.Code)

	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("錯誤的 JSON 請求:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
		return
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

// func Auth(username string, password string) error {
// 	if isExist := CheckUserIsExist(username); isExist {
// 		return CheckPassword(UserData[username], password)
// 	} else {
// 		return errors.New("user is not exist")
// 	}
// }

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
			"error": errors.New("必須輸入密碼"),
		})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if lockoutEnd, isLocked := lockoutTime[username]; isLocked {
		if time.Now().Before(lockoutEnd) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "帳號已被鎖定，請稍後再試",
			})
			return
		}
		// 重置鎖定時間和嘗試次數
		delete(lockoutTime, username)
		loginAttempts[username] = 0
	}

	if err := Auth(username, password); err == nil {
		loginAttempts[username] = 0 // 成功登入後重置嘗試次數
		c.HTML(http.StatusOK, "login.html", gin.H{
			"success": "登入成功",
		})
		return
	} else {
		loginAttempts[username]++
		if loginAttempts[username] >= 5 {
			lockoutTime[username] = time.Now().Add(5 * time.Minute)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "身分認證失敗，帳號已被鎖定五分鐘",
			})
			return
		}
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err,
		})
		return
	}
}

func CheckUserIsExist(username string) bool {
	_, isExist := UserData[username]
	return isExist
}

func CheckPassword(p1 string, p2 string) error {
	if p1 == p2 {
		return nil
	} else {
		return errors.New("密碼錯誤")
	}
}

func main() {

	InitDB("./database.db")
	defer DB.Close()

	fmt.Println("Data inserted successfully!")

	server := gin.Default()
	server.LoadHTMLGlob("./login.html")
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.POST("/resend-code", resendCodeHandler)
	server.POST("/send-code", sendCodeHandler)
	server.POST("/verify-code", verifyCodeHandler)
	server.Run(":8888")
}
