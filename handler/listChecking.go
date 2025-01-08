package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
)

func ListChecking(c *gin.Context) {
	// 驗證身份
	role, exist := c.Get("permission")
	if !exist || (role != "2" && role != "4") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 獲取 lecture_id
	lectureID := c.Param("lecture_id")
	if lectureID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing lecture ID"})
		return
	}

	// 查詢資料庫
	query := `
		SELECT u.user_id, u.user_name, ul.sign_in_time, ul.sign_out_time
		FROM Users_Lectures ul
		JOIN Users u ON ul.user = u.user_id
		WHERE ul.lecture = ?`
	rows, err := db.DB.Query(query, lectureID)
	if err != nil {
		fmt.Printf("Error querying database: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	defer rows.Close()

	// 解析查詢結果
	var records []map[string]interface{}
	for rows.Next() {
		var userID, userName string
		var signInTime, signOutTime sql.NullString

		err := rows.Scan(&userID, &userName, &signInTime, &signOutTime)
		if err != nil {
			fmt.Printf("Error scanning row: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}

		// 添加到結果集
		records = append(records, map[string]interface{}{
			"user_id":       userID,
			"user_name":     userName,
			"sign_in_time":  signInTime.String,
			"sign_out_time": signOutTime.String,
		})
	}

	// 返回結果
	c.JSON(http.StatusOK, gin.H{
		"lecture_id": lectureID,
		"records":    records,
	})
}
