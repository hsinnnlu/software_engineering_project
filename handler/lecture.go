package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/service"
)

func Lecturehandler(c *gin.Context) {
	// 從 URL 獲取參數
	lectureID := c.Param("lecture_id")
	userID := c.Param("user_id")

	// 從 JSON 輸入中提取其他內容
	var input struct {
		Sign_in_time  string `json:"sign_in_time"`
		Sign_out_time string `json:"sign_out_time"`
		Status        string `json:"status"`
	}

	// 綁定 JSON 資料
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 打印提取的參數
	fmt.Printf("Lecture ID: %s, User ID: %s, Status: %s, Sign-in Time: %s\n", lectureID, userID, input.Status, input.Sign_in_time)

	// 根據業務邏輯處理
	if input.Status == "in" {
		// 執行簽到邏輯
		err := service.InsertStudentIn(userID, lectureID, input.Sign_in_time, input.Sign_out_time)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign in"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Signed in successfully"})
		return
	}

	if input.Status == "out" {
		// 執行簽退邏輯
		err := service.InsertStudentIn(userID, lectureID, input.Sign_out_time, input.Sign_in_time)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign out"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
		return
	}

	// 處理無效的狀態
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
}
