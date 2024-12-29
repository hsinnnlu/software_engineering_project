package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/auth"
	"github.com/hsinnnlu/software_engineering_project/service"
)

func ResetpasswdSendlinkHandler(c *gin.Context) {
	var register_input struct {
		Email string `json:"email"`
	}

	// 前端的資料是否有空值
	if err := c.BindJSON(&register_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 檢查 email 有沒有存在資料庫中
	oldpassword, authenticated, err := auth.AuthenticateEmail(register_input.Email)
	if err != nil || !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成token
	token, err := service.GenerateResetpasswdToken(c, register_input.Email, oldpassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password token generation error"})
		return
	}

	// 傳送修改密碼連結
	err = service.SendResetLink(register_input.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "send mail error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重設密碼的連結已發送"})
}

func ResetpasswdVerifylinkHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid link"})
		return
	}

	// 驗證token
	claims, err := service.VerifylinkToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token verified successfully",
		"data":    claims, // 將 claims 作為回應數據
	})
}

func ResetpasswdChangepasswd(c *gin.Context) {
	var changepasswd_input struct {
		Password    string `json:"password"`
		Compassword string `json:"compassword"`
	}

	// 前端的資料是否有空值
	if err := c.BindJSON(&changepasswd_input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

}
