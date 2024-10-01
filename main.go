package main

import (
	"fmt"
	"net/http"

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

	// 設定靜態文件路由
	router.Static("/style.css", "./style.css")
	router.Static("/picture", "./picture")

	// 加載 HTML 模板
	router.LoadHTMLGlob("./webpage/**/*")

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)

	// 針對學生的路由
	router.GET("/webpage/Student/student.html", service.RoleMiddleware("1"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "student.html", nil)
	})

	// 針對管理員的路由
	router.GET("/webpage/manager/manager.html", service.RoleMiddleware("2"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "manager.html", nil)
	})

	// 針對教授的路由
	router.GET("/webpage/Professer/professer.html", service.RoleMiddleware("3"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "professer.html", nil)
	})

	router.GET("/student", service.StudentPage)

	// 測試部分
	router.GET("/test", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "Attendance_record.html", nil)
	})

	router.POST("/addLecture", func(c *gin.Context) {
		fmt.Println("call addLecture")

		lecture := models.Lecture{
			Id:        "002",
			Name:      "軟體工程",
			Timestamp: "星期三 13:00~16:00", // 時間格式還沒確定
			Location:  "資電館",
			Speaker:   "劉信宏",
		}
		db.InsertLecture(c, lecture)
	})

	// 將通配符路由放在最後
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil) // 可以設置404頁面
	})

	router.Run(":8080")
}