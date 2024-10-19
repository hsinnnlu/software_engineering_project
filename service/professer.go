package service

import (

	"github.com/gin-gonic/gin"
)

func Student_Attendance_record(c *gin.Context) {
	c.HTML(200, "Student_Attendance_record.html", nil)
}

func ShowStudentAttendanceRecord(c *gin.Context) {
	c.HTML(200, "Student_Attendance_record.html", nil)
}