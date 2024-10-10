package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"

	_ "github.com/mattn/go-sqlite3"
)

// 顯示個人檔案還沒做 2024.10.10

// header更新密碼
func HeaderResetPassword(c *gin.Context) {
	DB = db.DB
	session := sessions.Default(c)
	fmt.Println("pass HeaderResetPassword")

	// 如果 session.Get("user_id") 返回的值並不是 string 類型，這個類型斷言會引發 panic。 //2024.10.04 ???
	// 從 session 中獲取 user_id，進行類型斷言並檢查
	user_id, ok := session.Get("user_id").(string)
	user, err := db.GetUserById(DB, user_id)
	log.Println("user_id: ", user_id)
	log.Println("user_id_log: ", ok)
	log.Println("user_err", err)
	if !ok {
		c.HTML(http.StatusInternalServerError, "student.html", gin.H{
			"error": "無效的 user_id",
		})
		return
	}

	// 從前端請求獲取更改密碼的資料
	requestBody, err := HeaderResetPasswordRequest(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "hstudent.html", gin.H{
			"error": err,
			"user":  user,
		})
		return
	}

	// 密碼驗證
	err = validatePasswordMatch(requestBody.Password, requestBody.ConfirmPassword)
	log.Println("new_old_password_err: ", err)
	if err != nil {
		c.HTML(http.StatusBadRequest, "student.html", gin.H{
			"error": err,
			"user":  user,
		})
		return
	}

	// hash password
	hashedPassword := GetHashedPassword(requestBody.Password)

	// 更新數據庫中的密碼
	err = db.UpdatePasswordByUserid(DB, user_id, hashedPassword)
	if err != nil {
		c.HTML(http.StatusBadRequest, "student.html", gin.H{
			"error": err,
			"user":  user,
		})
	}

	// 清除 session
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	
	c.HTML(http.StatusOK, "student.html", gin.H{
		"message": "密碼重設成功，請重新登入",
	})
}

// 登出功能
func Logout(c *gin.Context) {
	// 清除 session
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	c.HTML(http.StatusOK, "student.html", gin.H{
		"message": "您已成功登出",
	})
}

func HeaderResetPasswordRequest(c *gin.Context) (struct {
	Password        string `form:"new-password"`
	ConfirmPassword string `form:"confirm-password"`
}, error) {
	var requestBody struct {
		Password        string `form:"new-password"`
		ConfirmPassword string `form:"confirm-password"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return requestBody, errors.New("無效的請求")
	}

	return requestBody, nil
}
