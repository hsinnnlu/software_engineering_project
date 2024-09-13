// package main

// import (
// 	"crypto/sha256"
// 	"database/sql"
// 	"encoding/hex"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"regexp"
// 	"sync"
// 	"time"

// 	"github.com/gin-contrib/sessions"
// 	"github.com/gin-contrib/sessions/cookie"

// 	"github.com/dgrijalva/jwt-go"

// 	"github.com/gin-gonic/gin"
// 	"gopkg.in/gomail.v2"

// 	_ "github.com/mattn/go-sqlite3"
// )

// type VerificationCode struct {
// 	Code       string
// 	Expiration time.Time
// }

// var Session = make(map[string]string) // session_id -> user_id

// var UserData map[string]string
// var verifyCodes map[string]VerificationCode

// type User struct {
// 	user_id    string
// 	password   string
// 	email      string
// 	permission string
// }

// var loginAttempts map[string]int
// var lockoutTime map[string]time.Time
// var mu sync.Mutex
// var jwtSecret = "f321a343233594d17697e0c9b83b6cb192a00e8562e4b1738e263c6ac90d3d1d"
// var DB *sql.DB

// func checkPassword(inputPassword, storedPasswordHash string) error {
// 	hashedInputPassword := getHashedPassword(inputPassword)
// 	fmt.Printf(hashedInputPassword)
// 	if hashedInputPassword != storedPasswordHash {
// 		return errors.New("password is incorrect")
// 	}
// 	return nil
// }

// func init() {
// 	verifyCodes = make(map[string]VerificationCode)
// 	loginAttempts = make(map[string]int)
// 	lockoutTime = make(map[string]time.Time)

// 	// 連接資料庫
// 	DB, err := sql.Open("sqlite3", "../database.db")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()

// }

// // 用戶驗證
// func Auth(user_id string, inputPassword string) (User, error) {
// 	var user User
// 	err := DB.QueryRow("SELECT user_id, password, user_permission FROM Users WHERE user_id = ?", user_id).Scan(&user.user_id, &user.password, &user.permission)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return user, errors.New("user does not exist")
// 		}
// 		return user, err
// 	}

// 	// 使用 checkPassword 函數來驗證密碼
// 	err = checkPassword(inputPassword, user.password)
// 	if err != nil {
// 		return user, err
// 	}

// 	return user, nil
// }

// // 登入頁面
// func LoginPage(c *gin.Context) {
// 	c.HTML(http.StatusOK, "login.html", nil)
// }

// func RoleMiddleware(requiredPermission string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session := sessions.Default(c)
// 		role := session.Get("role")

// 		if role == nil || role != requiredPermission {
// 			c.HTML(http.StatusForbidden, "login.html", gin.H{
// 				"error": "您沒有權限訪問該頁面",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

// // 登入驗證
// func LoginAuth(c *gin.Context) {
// 	var (
// 		user_id  string
// 		password string
// 	)
// 	if in, isExist := c.GetPostForm("user_id"); isExist && in != "" {
// 		user_id = in
// 	} else {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": errors.New("必須輸入使用者名稱"),
// 		})
// 		return
// 	}
// 	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
// 		password = in
// 	} else {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": errors.New("必須輸入密碼"),
// 		})
// 		return
// 	}

// 	mu.Lock()
// 	defer mu.Unlock()

// 	if lockoutEnd, isLocked := lockoutTime[user_id]; isLocked {
// 		if time.Now().Before(lockoutEnd) {
// 			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
// 				"error": "帳號已被鎖定，請稍後再試",
// 			})
// 			return
// 		}
// 		// 重置鎖定時間和嘗試次數
// 		delete(lockoutTime, user_id)
// 		loginAttempts[user_id] = 0
// 	}

// 	user, err := Auth(user_id, password)
// 	if err == nil {
// 		// 成功登入後重置嘗試次數
// 		loginAttempts[user_id] = 0

// 		session := sessions.Default(c)
// 		session.Set("user_id", user_id)
// 		session.Set("role", user.permission) // 保存用戶角色
// 		session.Save()

// 		// 根據角色進行跳轉
// 		switch user.permission {
// 		case "1":
// 			c.Redirect(http.StatusFound, "./webpage/Student/student.html")
// 		case "2":
// 			c.Redirect(http.StatusFound, "./webpage/manager/Account_manage.html")
// 		case "3":
// 			c.Redirect(http.StatusFound, "./webpage/Professer/Student_Attendance_record.html")
// 		default:
// 			c.HTML(http.StatusForbidden, "login.html", gin.H{
// 				"error": "未知的使用者角色",
// 			})
// 		}
// 		return
// 	} else {
// 		loginAttempts[user_id]++
// 		if loginAttempts[user_id] >= 5 {
// 			lockoutTime[user_id] = time.Now().Add(5 * time.Minute)
// 			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
// 				"error": "身分認證失敗，帳號已被鎖定五分鐘",
// 			})
// 			return
// 		}
// 		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
// 			"error": err,
// 		})
// 		return
// 	}
// }

// // 將密碼使用 SHA256
// func getHashedPassword(password string) string {
// 	hash := sha256.New()
// 	hash.Write([]byte(password))
// 	return hex.EncodeToString(hash.Sum(nil))
// }

// // 生成 JWT token，使用明文密碼
// func generateToken(email, password string) (string, error) {
// 	secretKey := jwtSecret + password // 使用原始密碼生成密鑰
// 	payload := jwt.MapClaims{
// 		"email":       email,
// 		"expiredTime": time.Now().Add(time.Hour).Unix(),
// 		"type":        "ResetPassword",
// 	}

// 	// 創建 Token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	return token.SignedString([]byte(secretKey))
// }

// func verifyToken(tokenString string) (jwt.MapClaims, error) {
// 	secretKeyFunc := func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("unexpected signing method")
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		email := claims["email"].(string)

// 		// 查詢資料庫獲取原始密碼
// 		var password string
// 		err := DB.QueryRow("SELECT password FROM Users WHERE user_email = ?", email).Scan(&password)
// 		if err != nil {
// 			return nil, errors.New("無法查詢用戶資料")
// 		}

// 		secretKey := jwtSecret + password

// 		return []byte(secretKey), nil
// 	}

// 	token, err := jwt.Parse(tokenString, secretKeyFunc)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if exp, ok := claims["expiredTime"].(float64); ok {
// 			expTime := time.Unix(int64(exp), 0)
// 			if time.Now().After(expTime) {
// 				return nil, errors.New("token has expired")
// 			}
// 		}
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token")
// }

// // 發送郵件
// func sendMail(to, body string) error {
// 	from := "nmg@cs.thu.edu.tw"
// 	password := "e04su3su;6"
// 	subject := "驗證碼"

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", to)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/plain", body)

// 	d := gomail.NewDialer("mail.cs.thu.edu.tw", 587, from, password)

// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // 發送重設密碼連結
// func sendResetLinkHandler(c *gin.Context) {
// 	var requestBody struct {
// 		Email string `json:"email"`
// 	}

// 	if err := c.BindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
// 		return
// 	}

// 	// 查詢資料庫獲取原始密碼
// 	var user User
// 	err := DB.QueryRow("SELECT user_id, password FROM Users WHERE user_email = ?", requestBody.Email).Scan(&user.user_id, &user.password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用戶不存在"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查詢資料庫失敗"})
// 		return
// 	}

// 	// 生成 JWT Token
// 	tokenString, err := generateToken(requestBody.Email, user.password)
// 	if err != nil {
// 		log.Println("生成 Token 失敗:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成驗證 Token 失敗"})
// 		return
// 	}

// 	// 生成重設密碼的連結
// 	resetLink := fmt.Sprintf("http://localhost:8888/login?token=%s", tokenString)

// 	// 發送郵件
// 	body := fmt.Sprintf("請點擊以下連結來重設您的密碼: %s", resetLink)

// 	err = sendMail(requestBody.Email, body)
// 	if err != nil {
// 		log.Println("發送郵件失敗:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "發送郵件失敗"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "重設密碼的連結已發送"})
// }

// // 重設密碼頁面
// func ResetPasswordPage(c *gin.Context) {
// 	token := c.Query("token")
// 	if token == "" {
// 		c.HTML(http.StatusBadRequest, "error.html", gin.H{
// 			"error": "無效的重設連結",
// 		})
// 		return
// 	}

// 	// 驗證Token
// 	claims, err := verifyToken(token)
// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "error.html", gin.H{
// 			"error": "Token 驗證失敗或已過期",
// 		})
// 		return
// 	}

// 	// Token有效，顯示更改密碼頁面
// 	c.HTML(http.StatusOK, "reset_password.html", gin.H{
// 		"token": token,           // 在HTML表單中以隱藏字段的形式保存token
// 		"email": claims["email"], // 顯示Email或用於後續驗證
// 	})
// }

// // 驗證密碼
// func validatePassword(password string) error {
// 	// 密碼長度至少8個字元
// 	if len(password) < 8 {
// 		log.Println("Password:", password)
// 		return errors.New("密碼長度必須至少8個字元")
// 	}

// 	// 密碼必須包含至少一個大寫字母
// 	uppercasePattern := `[A-Z]`
// 	matched, _ := regexp.MatchString(uppercasePattern, password)
// 	if !matched {
// 		return errors.New("密碼必須包含至少一個大寫字母")
// 	}

// 	// 密碼不可包含空白字元或特殊符號
// 	illegalPattern := `[^\w]`
// 	matched, _ = regexp.MatchString(illegalPattern, password)
// 	if matched {
// 		return errors.New("密碼不可包含空白字元或特殊符號")
// 	}

// 	return nil
// }

// func ResetPasswordHandler(c *gin.Context) {
// 	var requestBody struct {
// 		Token           string `form:"token"`
// 		Email           string `form:"email"`
// 		Password        string `form:"password"`
// 		ConfirmPassword string `form:"confirm-password"`
// 	}

// 	if err := c.Bind(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的請求"})
// 		return
// 	}

// 	if err := validatePassword(requestBody.Password); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
// 		return
// 	}

// 	if requestBody.Password != requestBody.ConfirmPassword {
// 		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "密碼和確認密碼不一致"})
// 		return
// 	}

// 	// 解碼 Token 並提取 Email
// 	claims, err := verifyToken(requestBody.Token)
// 	if err != nil {
// 		log.Println("Token 驗證失敗:", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "無效的驗證 Token"})
// 		return
// 	}

// 	email, ok := claims["email"].(string)
// 	if !ok {
// 		log.Println("從 Token 中提取 Email 失敗")
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Token 資料無效"})
// 		return
// 	}

// 	hashedPassword := getHashedPassword(requestBody.Password)

// 	// 更新密碼
// 	result, err := DB.Exec("UPDATE Users SET password = ? WHERE user_email = ?", hashedPassword, email)
// 	if err != nil {
// 		log.Println("更新密碼失敗:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新密碼失敗"})
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		log.Println("無法獲取受影響的行數:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "密碼更新操作失敗"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		log.Println("更新密碼失敗: 無法找到對應的 email 或者密碼未修改")
// 		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "無法找到對應的 email"})
// 		return
// 	}

// 	log.Println("密碼更新成功")
// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "密碼已成功重設"})
// }

// // 主程式
// func main() {

// 	server := gin.Default()
// 	// 使用 Gin 的 session
// 	store := cookie.NewStore([]byte("secret"))
// 	server.Use(sessions.Sessions("mysession", store))

// 	server.StaticFile("/style.css", "./style.css")
// 	server.Static("/picture", "./picture")

// 	server.LoadHTMLFiles(
// 		"./login.html",
// 		"./reset_password.html",
// 	)

// 	server.LoadHTMLGlob(
// 		"./webpage/**/*",
// 	)

// 	server.GET("/login", LoginPage)
// 	server.POST("/login", LoginAuth)

// 	server.POST("/resend-code", sendResetLinkHandler)
// 	server.POST("/send-code", sendResetLinkHandler)
// 	server.GET("/reset-password", ResetPasswordPage)
// 	server.POST("/reset-password", ResetPasswordHandler)

// 	// 學生頁面
// 	server.GET("/webpage/Student/student.html", RoleMiddleware("1"), func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "student.html", nil)
// 	})

// 	// 管理員頁面
// 	server.GET("/webpage/manager/Account_manage.html", RoleMiddleware("2"), func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "Account_manage.html", nil)
// 	})

// 	// 教授頁面
// 	server.GET("/webpage/Professer/Student_Attendance_record.html", RoleMiddleware("3"), func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "Student_Attendance_record.html", nil)
// 	})
// 	server.Run(":8888")
// }
