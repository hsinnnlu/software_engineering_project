// POST: /login 路由
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/models"
)

// 登入頁面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginAuth(c *gin.Context) {

	// 輸入邏輯處理: 是否輸入使用者的帳號密碼
	input_id, input_password := preProcessingInput(c)
	fmt.Println("pass A")

	// 檢查密碼是否正確
	user := &models.User{}
	user, err := checkPassword(input_id, input_password)
	fmt.Print()

	// 錯誤處理
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err,
		})
		return
		// 嘗試次數計數：還沒做
	}

	// 登入成功: 設定 session

	session := sessions.Default(c)
	session.Set("user_id", user.Id)      // 保存用戶名稱
	session.Set("role", user.Permission) // 保存用戶角色
	session.Save()

}

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
	hashedInputPassword := getHashedPassword(inputPassword)
	fmt.Println(hashedInputPassword) // 這裡會印出輸入密碼的 SHA256 雜湊值

	// 檢查使用者是否存在
	user := &models.User{}
	user, err := db.GetUserById(user_id)
	if err != nil {
		return nil, errors.New("user does not exist")
	}

	storedPasswordHash := user.Password_hash
	if storedPasswordHash != hashedInputPassword {
		return nil, errors.New("password is incorrect")
	}
	return user, nil
}

// 將密碼使用 SHA256
func getHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
