package main

import (
	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/software_engineering_project/auth"
	"github.com/hsinnnlu/software_engineering_project/db"

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
	//router.POST("/test", service.TestDB)

	router.Run(":8080")
}
