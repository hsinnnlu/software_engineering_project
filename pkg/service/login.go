// POST: /login 路由

package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var loginAttempts = make(map[string]int)
var lockoutTime = make(map[string]time.Time)

var mu sync.Mutex // 互斥鎖

func getHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func checkPassword(inputPassword, storedPasswordHash string) error {
	hashedInputPassword := getHashedPassword(inputPassword)
	fmt.Printf(hashedInputPassword)
	if hashedInputPassword != storedPasswordHash {
		return errors.New("password is incorrect")
	}
	return nil
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

	user, err := SelectUser(DB, user_id)
	if err == nil {
		// 成功登入後重置嘗試次數
		loginAttempts[user_id] = 0

		session := sessions.Default(c)
		session.Set("user_id", user_id)
		session.Set("role", user.permission) // 保存用戶角色
		session.Save()

		// 根據角色進行跳轉
		switch user.permission {
		case "1":
			c.Redirect(http.StatusFound, "./webpage/Student/student.html")
		case "2":
			c.Redirect(http.StatusFound, "./webpage/manager/Account_manage.html")
		case "3":
			c.Redirect(http.StatusFound, "./webpage/Professer/Student_Attendance_record.html")
		default:
			c.HTML(http.StatusForbidden, "login.html", gin.H{
				"error": "未知的使用者角色",
			})
		}
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
