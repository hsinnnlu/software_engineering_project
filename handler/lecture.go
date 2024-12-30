package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/service"
)

func Lecturehandler(c *gin.Context) {
	var input struct {
		User_id       string `json:"user_id"`
		Lecture_id    string `json:"lecture_id"`
		Sign_in_time  string `json:"sign_in_time"`
		Sign_out_time string `json:"sign_out_time"`
		Status        string `json:"status"`
	}

	// 前端的資料是否有空值
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := service.Authlecture(input.Lecture_id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid lecture_id"})
		return
	}
}
