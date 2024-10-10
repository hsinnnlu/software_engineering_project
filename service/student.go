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

func Lecture_information(c *gin.Context) {
	session := sessions.Default(c)
	redirectURL, _ := session.Get("redirect_url").(string)
	
	c.HTML(200, "Lecture_information.html", gin.H{
        "redirect_url": redirectURL,
    })
}

func Lecture_notes(c *gin.Context) {
	c.HTML(200, "Lecture_notes.html", nil)
}