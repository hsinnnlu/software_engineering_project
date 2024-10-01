// POST: /login 路由
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/models"
)

var DB = db.DB

// 登入頁面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginAuth(c *gin.Context) {
	// 初始化 session
	session := sessions.Default(c) // 在這裡直接取得 session

	// 輸入邏輯處理: 是否輸入使用者的帳號密碼
	input_id, input_password, err := preProcessingInput(c)
	if err != nil {
		return
	}

	// 檢查帳號是否已鎖定
	if checkLockoutStatus(c) { // 傳遞 gin.Context
		c.HTML(http.StatusTooManyRequests, "login.html", gin.H{
			"error": "帳號已鎖定，請稍後再試。",
		})
		return
	}

	locktime := session.Get("locktime")
	log.Print("locktimetest:", locktime)

	// 驗證密碼
	user, err := checkPassword(input_id, input_password)
	if err != nil {
		// 增加登入失敗次數
		if incrementLoginAttempts(c) { // 傳遞 gin.Context
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "帳號已被鎖定三分鐘，請稍後再試。",
			})
		} else {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "帳號或密碼錯誤，請再試一次。",
			})
		}
		return
	}

	// 登入成功: 重置計數器和 session 狀態
	c.HTML(http.StatusOK, "success.html", gin.H{
		"message": "登入成功",
	})

	session.Set("login_attempts", 0)
	session.Delete("locktime")
	session.Set("user_id", user.Id)
	session.Set("role", user.Permission)
	session.Save()

	// 根據權限進行重定向
	RedirectByPermission(c, *user)
}

// 測試成功 2024/09/24
// 資料欄位預處理：檢查是否有輸入帳號密碼
func preProcessingInput(c *gin.Context) (user_id, password string, err error) {
	if in, isExist := c.GetPostForm("user_id"); isExist && in != "" {
		user_id = in
	} else {
		err = errors.New("必須輸入使用者名稱")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err,
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		err = errors.New("必須輸入密碼")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err,
		})
		return
	}
	return user_id, password, nil
}

// 測試成功 2024/09/24
// 將密碼使用 SHA256
func GetHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// 測試成功 2024/09/24
func checkPassword(user_id, inputPassword string) (*models.User, error) {
	hashedInputPassword := GetHashedPassword(inputPassword)
	fmt.Println(hashedInputPassword) // 這裡會印出輸入密碼的 SHA256 雜湊值

	// 檢查使用者是否存在
	user, err := db.GetUserById(DB, user_id)
	if err != nil {
		fmt.Println("error:", err)
		return nil, errors.New("user does not exist")
	}
	fmt.Println("pass C: ", user.Password_hash) // 這裡會印出資料庫中的密碼雜湊值
	storedPasswordHash := user.Password_hash
	if storedPasswordHash != hashedInputPassword {
		return nil, errors.New("password is incorrect")
	}
	return user, nil
}

// senssion沒有time的資料
// 檢查是否在鎖定時間內
func checkLockoutStatus(c *gin.Context) bool {
	session := sessions.Default(c)
	locktime := session.Get("locktime")
	log.Print("locktime:", locktime)

	if locktime != nil {
		lockoutEndTime := locktime.(time.Time)
		if time.Now().Before(lockoutEndTime) {
			return true
		} else {
			session.Set("login_attempts", 0)
			session.Delete("locktime")
			session.Save()
		}
	}
	return false
}

// 登入失敗計數器
func incrementLoginAttempts(c *gin.Context) bool {
	session := sessions.Default(c)
	attempts := session.Get("login_attempts")

	var attemptsCount int
	if attempts == nil {
		attemptsCount = 0
	} else {
		attemptsCount = attempts.(int)
	}

	attemptsCount++
	session.Set("login_attempts", attemptsCount)
	log.Print("login_attempts 設置為: ", session.Get("login_attempts"))

	if attemptsCount >= 5 {
		session.Set("locktime", time.Now().Add(3*time.Minute))
		log.Print("locktime 設置為: ", session.Get("locktime"))
		session.Save()
		return true
	}

	session.Save()
	return false
}

// 根據身份進行重導向
func RedirectByPermission(c *gin.Context, user models.User) {

	// userInfo := map[string]string{
	// 	"user_id":    user.Id,
	// 	"permission": user.Permission,
	// 	"name":       user.Name,
	// }

	switch user.Permission {
	case "1":
		c.Redirect(http.StatusFound, "/webpage/Student/student.html")
		// c.HTML(http.StatusFound, "student.html", gin.H{
		// 	"user": userInfo,
		// })
	case "2":
		c.Redirect(http.StatusFound, "/webpage/manager/manager.html")
	case "3":
		c.Redirect(http.StatusFound, "/webpage/Professer/professer.html")
	default:
		c.HTML(http.StatusForbidden, "login.html", gin.H{
			"error": "未知的使用者角色",
		})
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")

		if role == nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "您尚未登入，請先登入。",
			})
			c.Abort()
			return
		}

		userRole := role.(string)
		allowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.HTML(http.StatusForbidden, "switch_url_error.html", gin.H{
				"error": "您沒有權限訪問該頁面。",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
