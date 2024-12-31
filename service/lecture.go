package service

import (
	"database/sql"
	"fmt"

	"github.com/hsinnnlu/software_engineering_project/db"
)

func Authlecture(lecture_id string) error {
	var lecture string

	query := `SELECT lecture_name FROM Lectures WHERE lecture_id = ?`
	err := db.DB.QueryRow(query, lecture_id).Scan(&lecture)

	if err == sql.ErrNoRows {
		fmt.Printf("err: %s\n", err)
		return nil // 用戶不存在
	} else if err != nil {
		fmt.Printf("err: %s\n", err)
		return err // 一些奇怪的錯誤報錯
	}

	return nil
}

func GetActiveLectures() ([]struct{ ID, Name string }, error) {
	var lectures []struct{ ID, Name string }

	query := `SELECT lecture_id, lecture_name FROM Lectures WHERE status = 1`
	rows, err := db.DB.Query(query)
	if err != nil {
		fmt.Printf("Query error: %s\n", err)
		return nil, err // 回傳錯誤
	}
	defer rows.Close()

	for rows.Next() {
		var lecture struct{ ID, Name string }
		if err := rows.Scan(&lecture.ID, &lecture.Name); err != nil {
			fmt.Printf("Row scan error: %s\n", err)
			return nil, err // 回傳錯誤
		}
		lectures = append(lectures, lecture)
	}

	// 檢查 rows 是否發生錯誤
	if err := rows.Err(); err != nil {
		fmt.Printf("Rows error: %s\n", err)
		return nil, err
	}

	return lectures, nil
}

func InsertStudentIn(userID, lectureID, signInTime string) error {
	// 檢查記錄是否已存在
	queryCheck := `SELECT COUNT(*) FROM Users_Lectures WHERE user = ? AND lecture = ?`
	var count int
	err := db.DB.QueryRow(queryCheck, userID, lectureID).Scan(&count)
	if err != nil {
		fmt.Printf("Check error: %s\n", err)
		return err
	}

	if count > 0 {
		// 更新簽到時間
		queryUpdate := `UPDATE Users_Lectures SET sign_in_time = ? WHERE user = ? AND lecture = ?`
		_, err = db.DB.Exec(queryUpdate, signInTime, userID, lectureID)
		if err != nil {
			fmt.Printf("Update error: %s\n", err)
			return err
		}
		return nil
	}

	// 插入新記錄
	queryInsert := `INSERT INTO Users_Lectures (user, lecture, sign_in_time) VALUES (?, ?, ?)`
	_, err = db.DB.Exec(queryInsert, userID, lectureID, signInTime)
	if err != nil {
		fmt.Printf("Insert error: %s\n", err)
		return err
	}

	return nil
}

func InsertStudentOut(userID, lectureID, signOutTime string) error {
	// 檢查記錄是否已存在
	queryCheck := `SELECT COUNT(*) FROM Users_Lectures WHERE user = ? AND lecture = ?`
	var count int
	err := db.DB.QueryRow(queryCheck, userID, lectureID).Scan(&count)
	if err != nil {
		fmt.Printf("Check error: %s\n", err)
		return err
	}

	if count > 0 {
		// 更新簽退時間
		queryUpdate := `UPDATE Users_Lectures SET sign_out_time = ? WHERE user = ? AND lecture = ?`
		_, err = db.DB.Exec(queryUpdate, signOutTime, userID, lectureID)
		if err != nil {
			fmt.Printf("Update error: %s\n", err)
			return err
		}
		return nil
	}

	// 插入新記錄
	queryInsert := `INSERT INTO Users_Lectures (user, lecture, sign_out_time) VALUES (?, ?, ?)`
	_, err = db.DB.Exec(queryInsert, userID, lectureID, signOutTime)
	if err != nil {
		fmt.Printf("Insert error: %s\n", err)
		return err
	}

	return nil
}

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hsinnnlu/software_engineering_project/db"
// 	"github.com/hsinnnlu/software_engineering_project/models"
// )

// // 新增講座：目前只有 名稱、時間

// // func LectureList() []models.Lecture {

// // }

// func AddLecture(c *gin.Context) {
// 	fmt.Print("AddLecture\n")

// 	lecture := models.Lecture{}
// 	// 講座編號是自動產生的，不需要輸入

// 	if in, isExist := c.GetPostForm("lecture_name"); isExist && in != "" {
// 		lecture.Name = c.PostForm("lecture_name")
// 	} else {
// 		c.HTML(http.StatusBadRequest, "manager.html", gin.H{
// 			"error": "必須輸入講座名稱",
// 		})
// 		return
// 	}

// 	if in, isExist := c.GetPostForm("lecture_timestamp"); isExist && in != "" {
// 		lecture.Timestamp = c.PostForm("lecture_timestamp")
// 	} else {
// 		c.HTML(http.StatusBadRequest, "manager.html", gin.H{
// 			"error": "必須輸入講座時間",
// 		})
// 		return
// 	}
// 	err := db.InsertLecture(c, lecture)

// 	if err != nil {
// 		c.HTML(http.StatusBadRequest, "manager.html", gin.H{
// 			"error": err,
// 		})
// 		return
// 	}

// 	c.HTML(http.StatusCreated, "Lecture_manage.html", nil)
// }

// func ListLectures(c *gin.Context) {
// 	lectures, err := db.GetLecturesByStatus(DB, "0", "1", "2")

// 	if err != nil {
// 		fmt.Print("error!\n")
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	c.JSON(http.StatusFound, lectures)
// }
