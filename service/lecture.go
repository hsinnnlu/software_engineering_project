package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/models"
)

// 新增講座：目前只有 名稱、時間

// func LectureList() []models.Lecture {

// }

func AddLecture(c *gin.Context) {
	fmt.Print("AddLecture\n")

	lecture := models.Lecture{}
	// 講座編號是自動產生的，不需要輸入

	if in, isExist := c.GetPostForm("lecture_name"); isExist && in != "" {
		lecture.Name = c.PostForm("lecture_name")
	} else {
		c.HTML(http.StatusBadRequest, "manager.html", gin.H{
			"error": "必須輸入講座名稱",
		})
		return
	}

	if in, isExist := c.GetPostForm("lecture_timestamp"); isExist && in != "" {
		lecture.Timestamp = c.PostForm("lecture_timestamp")
	} else {
		c.HTML(http.StatusBadRequest, "manager.html", gin.H{
			"error": "必須輸入講座時間",
		})
		return
	}
	err := db.InsertLecture(c, lecture)

	if err != nil {
		c.HTML(http.StatusBadRequest, "manager.html", gin.H{
			"error": err,
		})
		return
	}

	c.HTML(http.StatusCreated, "Lecture_manage.html", nil)
}

func ListLectures(c *gin.Context) {
	lectures, err := db.GetLecturesByStatus(DB, "0", "1", "2")

	if err != nil {
		fmt.Print("error!\n")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusFound, lectures)
}
