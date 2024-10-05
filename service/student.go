package service

import (
	"fmt"

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
