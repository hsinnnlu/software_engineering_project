package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/software_engineering_project/auth"
	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/models"

	"github.com/hsinnnlu/software_engineering_project/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db.InitDB("./database.db")

	router := gin.Default()
	router.Use(auth.InitSession("secret"))

	router.Static("/style.css", "./style.css")
	router.Static("/picture", "./picture")
	router.LoadHTMLGlob(("./webpage/**/*"))

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)

	router.GET("/student", service.StudentPage)

	// 以下都是Test的部分
	router.GET("/test", func(ctx *gin.Context) {
		ctx.HTML(200, "Attendance_record.html", nil)
	})

	router.GET("/manager", func(c *gin.Context) {
		c.HTML(200, "manager.html", nil)
	})

	// 新增講座（目前用lambda，還沒包起來）

	router.POST("/ttt", service.AddLecture)
	router.POST("/addLecture", func(c *gin.Context) {

		fmt.Print("call addLecture\n")

		lecture := models.Lecture{
			Id:        3,
			Name:      "軟體工程",
			Timestamp: "星期三 13:00~16:00", // 時間格式還沒確定
			Location:  "資電館",
			Speaker:   "劉信宏",
		}
		db.InsertLecture(c, lecture)

		c.HTML(201, "manager.html", nil)
	})

	//router.POST("/test", service.TestDB)

	router.Run(":8080")
}
