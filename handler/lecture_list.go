package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/service"
)

// 列出有哪些講座
func Lecturelisthandler(c *gin.Context) {
	lectures, err := service.GetActiveLectures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lecture error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lecture": lectures})
}

// 列出參與過哪些講座
func LectureParticipatehandler(c *gin.Context) {
	type ParticipatedLecture struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		SignInTime  string `json:"sign_in_time"`
		SignOutTime string `json:"sign_out_time"`
		Place       string `json:"place"`
		Speaker     string `json:"speaker"`
	}
	// 確保身份是學生
	role, exist := c.Get("permission")
	if !exist || role != "1" { // 假設 "1" 是學生的權限代碼
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 獲取當前用戶 ID
	userID, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// 定義查詢 SQL
	query := `
		SELECT
			l.lecture_id, l.lecture_name,  
			COALESCE(ul.sign_in_time, '未簽到') AS sign_in_time, 
			COALESCE(ul.sign_out_time, '未簽退') AS sign_out_time,
			l.lecture_location,
			l.lecture_speaker
		FROM
			Users_Lectures ul
		JOIN
			Lectures l ON ul.lecture = l.lecture_id
		WHERE
			ul.user = ?
	`

	// 查詢參與的講座
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()

	// 解析數據
	var lectures []ParticipatedLecture
	for rows.Next() {
		var lecture ParticipatedLecture
		if err := rows.Scan(&lecture.ID, &lecture.Title, &lecture.SignInTime, &lecture.SignOutTime, &lecture.Place, &lecture.Speaker); err != nil {
			fmt.Printf("error: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data", "details": err.Error()})
			return
		}
		lectures = append(lectures, lecture)
	}

	// 檢查是否有任何參與記錄
	if len(lectures) == 0 {
		c.JSON(http.StatusOK, gin.H{"lectures": []ParticipatedLecture{}})
		return
	}

	// 返回結果
	c.JSON(http.StatusOK, gin.H{"lectures": lectures})
}
