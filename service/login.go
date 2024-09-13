// POST: /login 路由
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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

	// 輸入邏輯處理: 是否輸入使用者的帳號密碼
	input_id, input_password := preProcessingInput(c)
	fmt.Println("pass A")

	// 看帳號有沒有被鎖定，有的話就報錯
	isLocked := checkLockoutStatus(c)
	if isLocked {
		// 如果帳號已鎖定，顯示錯誤訊息
		c.HTML(http.StatusTooManyRequests, "login.html", gin.H{
			"error": "帳號已鎖定，請在三分鐘後重試。",
		})
		return
	}

	// 沒有的話再檢查密碼是否正確
	user := &models.User{} //這裡的變數用在哪裡 沒有用到哈哈，其實下一行user（type User)就是了
	user, err := checkPassword(input_id, input_password)
	fmt.Print()

	// 錯誤處理
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err,
		})

		// 計數器
		attempts := incrementLoginAttempts(c)
		if attempts {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "身分認證失敗，帳號已被鎖定三分鐘",
			})
		}
		return
	}

	// 登入成功: 設定 session
	session := sessions.Default(c)
	session.Set("user_id", user.Id)      // 保存用戶名稱
	session.Set("role", user.Permission) // 保存用戶角色
	session.Save()

	// 根據不同使用者跳至不同的介面
	RedirectByPermission(c, user.Permission)
}

// 資料欄位預處理：檢查是否有輸入帳號密碼
func preProcessingInput(c *gin.Context) (user_id, password string) {

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
	return user_id, password
}

// 有個資安小問題：待改進 2024/09/09
func checkPassword(user_id, inputPassword string) (*models.User, error) {
	hashedInputPassword := GetHashedPassword(inputPassword)
	fmt.Println(hashedInputPassword) // 這裡會印出輸入密碼的 SHA256 雜湊值

	// 檢查使用者是否存在
	user := &models.User{}
	user, err := db.GetUserById(DB, user_id)
	fmt.Println("pass B: ", user.Password_hash) // 這裡會印出資料庫中的密碼雜湊值
	if err != nil {
		fmt.Println("error:", err)
		return nil, errors.New("user does not exist")
	}

	storedPasswordHash := user.Password_hash
	if storedPasswordHash != hashedInputPassword {
		return nil, errors.New("password is incorrect")
	}
	return user, nil
}

// 根據身份進行重導向
func RedirectByPermission(c *gin.Context, userPermission string) {
	switch userPermission {
	case "1":
		c.Redirect(http.StatusFound, "/student")
	case "2":
		c.Redirect(http.StatusFound, "/webpage/manager/Account_manage.html")
	case "3":
		c.Redirect(http.StatusFound, "/webpage/Professer/Student_Attendance_record.html")
	default:
		c.HTML(http.StatusForbidden, "login.html", gin.H{
			"error": "未知的使用者角色",
		})
	}
}

// 檢查是否在鎖定時間內
func checkLockoutStatus(c *gin.Context) bool {
	session := sessions.Default(c)
	locktime := session.Get("locktime")

	// 檢查是否在鎖定時間內
	if locktime != nil {
		lockoutEndTime := locktime.(time.Time)
		if time.Now().Before(lockoutEndTime) {
			return true
		} else {
			// 如果鎖定時間已過，重置鎖定狀態
			session.Set("login_attempts", 0)
			session.Set("locktime", nil)
			session.Save()
		}
	}
	// 如果未鎖定，返回 false
	return false
}

// 將密碼使用 SHA256
func GetHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// 登入失敗計數器
func incrementLoginAttempts(c *gin.Context) bool {
	session := sessions.Default(c)
	attempts := session.Get("login_attempts")

	if attempts == nil {
		attempts = 0
	}

	attempts = attempts.(int) + 1
	session.Set("login_attempts", attempts)

	// 如果錯誤次數達到 5 次，鎖定帳號三分鐘
	if attempts.(int) >= 5 {
		session.Set("locktime", time.Now().Add(3*time.Minute))
		session.Save()
		return true
	}

	session.Save()
	return false
}
