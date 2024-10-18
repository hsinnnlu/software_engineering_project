package service

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
)

func RenderAnnouncement(c *gin.Context) {
	announceList, err := db.GetAnnouncementList(DB)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("announceList: %v\n", announceList)
	c.JSON(200, announceList)
}

func Attendance_record(c *gin.Context) {
	c.HTML(200, "Attendance_record.html", nil)
}

func ShowAttendanceRecord(c *gin.Context) {
	session := sessions.Default(c)
	user_id, _ := session.Get("user_id").(string)

	User_Lectures, err := db.GetAttendanceRecord(DB, user_id)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("User_Lectures: %v\n", User_Lectures)
	c.JSON(200, User_Lectures)
}

func Lecture_information(c *gin.Context) {
	session := sessions.Default(c)
	redirectURL, _ := session.Get("redirect_url").(string)

	c.HTML(200, "Lecture_information.html", gin.H{
		"redirect_url": redirectURL,
	})
}

func ShowLectureInformation(c *gin.Context) {
	session := sessions.Default(c)
	user_id, _ := session.Get("user_id").(string)

	User_Lectures, err := db.GetLectureInformation(DB, user_id)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("User_Lectures: %v\n", User_Lectures)
	c.JSON(200, User_Lectures)
}

func Lecture_notes(c *gin.Context) {
	c.HTML(200, "Lecture_notes.html", nil)
}
