package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
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
		fmt.Printf("error: %s\n", err)
		return
	}

	// 打印提取的參數
	fmt.Printf("Lecture ID: %s, User ID: %s, Status: %s, Sign-in Time: %s\n", lectureID, userID, input.Status, input.Sign_in_time)

	// 根據業務邏輯處理
	if input.Status == "in" {
		// 執行簽到邏輯
		err := service.InsertStudentIn(userID, lectureID, input.Sign_in_time)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign in"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Signed in successfully"})
		return
	}

	if input.Status == "out" {
		// 執行簽退邏輯
		err := service.InsertStudentOut(userID, lectureID, input.Sign_out_time)
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

// 修改表單
func EditLecture(c *gin.Context) {
	permission, exist := c.Get("permission")
	// 驗證身份
	if !exist || permission != "2" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		fmt.Printf("role: %s\n", permission)
		return
	}
	// 定義接收表單的結構
	type EditLectureRequest struct {
		ID    int    `json:"id" binding:"required"`    // 必須提供講座 ID
		Title string `json:"title" binding:"required"` // 講座名稱
		Date  string `json:"date" binding:"required"`  // 講座日期
		Time  string `json:"time" binding:"required"`  // 講座時間
		Place string `json:"place" binding:"required"` // 講座地點
	}

	// 接收 JSON 請求
	var req EditLectureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// UPDATE 資料庫
	query := `UPDATE Lectures SET lecture_name = ?, lecture_timestamp = ?, lecture_location = ? WHERE lecture_id = ?`
	res, err := db.DB.Exec(query, req.Title, req.Date+"T"+req.Time, req.Place, req.ID)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lecture", "details": err.Error()})
		return
	}

	// 檢查是否有更新到記錄
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch update result", "details": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lecture not found or no changes made"})
		return
	}

	// 成功更新
	c.JSON(http.StatusOK, gin.H{
		"message": "Lecture updated successfully",
		"lecture": req,
	})
}
