package service

import (

	"github.com/gin-gonic/gin"
)


func Lecture_CheckIn(c *gin.Context) {
	c.HTML(200, "Lecture_CheckIn.html", nil)
}

func Lecture_manage(c *gin.Context) {
	c.HTML(200, "Lecture_manage.html", nil)
}

func Account_manage(c *gin.Context) {
	c.HTML(200, "Account_manage.html", nil)
}