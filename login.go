package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"

	"database/sql"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	user_id  string
	password string
	email    string
}

var loginAttempts map[string]int
var lockoutTime map[string]time.Time
var mu sync.Mutex
var jwtSecret = "f321a343233594d17697e0c9b83b6cb192a00e8562e4b1738e263c6ac90d3d1d"
var DB *sql.DB

func init() {
	loginAttempts = make(map[string]int)
	lockoutTime = make(map[string]time.Time)
}

// 初始化資料庫
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
}

// 用戶驗證
func Auth(user_id string, password string) error {
	var user User
	err := DB.QueryRow("SELECT user_id, password FROM Users WHERE user_id = ?", user_id).Scan(&user.user_id, &user.password)
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

// 登入頁面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// 登入驗證
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

	if lockoutEnd, isLocked := lockoutTime[user_id]; isLocked {
		if time.Now().Before(lockoutEnd) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "帳號已被鎖定，請稍後再試",
			})
			return
		}
		// 重置鎖定時間和嘗試次數
		delete(lockoutTime, user_id)
		loginAttempts[user_id] = 0
	}

	if err := Auth(user_id, password); err == nil {
		loginAttempts[user_id] = 0 // 成功登入後重置嘗試次數
		c.HTML(http.StatusOK, "login.html", gin.H{
			"success": "登入成功",
		})
		return
	} else {
		loginAttempts[user_id]++
		if loginAttempts[user_id] >= 5 {
			lockoutTime[user_id] = time.Now().Add(5 * time.Minute)
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

// 將密碼使用 SHA256
func getHashedPassword(email string) string {
	hash := sha256.New()
	hash.Write([]byte(email))
	return hex.EncodeToString(hash.Sum(nil))
}

// 生成 JWT token
func generateToken(email string) (string, error) {
	userHashedPassword := getHashedPassword(email)
	secretKey := jwtSecret + userHashedPassword
	payload := jwt.MapClaims{
		"email":       email,
		"expiredTime": time.Now().Add(time.Hour).Unix(), // 設置過期時間為 1 小時後
		"type":        "ResetPassword",
	}

	// 創建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secretKey))
}

// 驗證 JWT token
func verifyToken(tokenString string) (jwt.MapClaims, error) {
	secretKeyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		userHashedPassword := getHashedPassword(email)
		secretKey := jwtSecret + userHashedPassword

		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(tokenString, secretKeyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["expiredTime"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if time.Now().After(expTime) {
				return nil, errors.New("token has expired")
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 發送郵件
func sendMail(to, body string) error {
	from := "nmg@cs.thu.edu.tw"
	password := "e04su3su;6" // 建議使用環境變數來存儲敏感信息
	subject := "驗證碼"

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

// 發送重設密碼連結
func sendResetLinkHandler(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
		return
	}

	// 生成 JWT Token
	tokenString, err := generateToken(requestBody.Email)
	if err != nil {
		log.Println("生成 Token 失敗:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成驗證 Token 失敗"})
		return
	}

	// 生成重設密碼的連結
	resetLink := fmt.Sprintf("http://localhost:8888/reset-password?token=%s", tokenString)

	// 發送郵件
	body := fmt.Sprintf("請點擊以下連結來重設您的密碼: %s", resetLink)

	err = sendMail(requestBody.Email, body)
	if err != nil {
		log.Println("發送郵件失敗:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "發送郵件失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "重設密碼的連結已發送"})
}

// 重設密碼頁面
func ResetPasswordPage(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "無效的重設連結",
		})
		return
	}

	// 顯示重設密碼的頁面
	c.HTML(http.StatusOK, "reset_password.html", gin.H{"token": token})
}

// 重設密碼處理
func ResetPasswordHandler(c *gin.Context) {
	var requestBody struct {
		Token    string `json:"token"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
		return
	}

	// 驗證 Token
	claims, err := verifyToken(requestBody.Token)
	if err != nil {
		log.Println("Token 驗證失敗:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "無效的驗證 Token"})
		return
	}

	// 驗證 Email 是否一致
	if claims["email"] != requestBody.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token 中的 Email 和請求中的 Email 不一致"})
		return
	}

	// 更新資料庫中的密碼
	hashedPassword := getHashedPassword(requestBody.Password)
	_, err = DB.Exec("UPDATE Users SET password = ? WHERE email = ?", hashedPassword, requestBody.Email)
	if err != nil {
		log.Println("更新密碼失敗:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新密碼失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "密碼已成功重設"})
}

// 主程式
func main() {
	InitDB("./database.db")
	defer DB.Close()

	server := gin.Default()
	server.StaticFile("/style.css", "./style.css")
	server.Static("/picture", "./picture")
	server.LoadHTMLFiles("./login.html", "./reset_password.html")
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.POST("/resend-code", sendResetLinkHandler)
	server.POST("/send-code", sendResetLinkHandler)
	server.GET("/reset-password", ResetPasswordPage)
	server.POST("/reset-password", ResetPasswordHandler)
	server.Run(":8888")
}
