package service

import "github.com/gin-gonic/gin"

func Logout(c *gin.Context) {
	// 清除 session
	c.SetCookie("session", "", -1, "/", "localhost", false, true)

}
