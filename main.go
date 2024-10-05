package main

import (
	"net/http"

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

	// 設定靜態文件路由
	router.Static("/style.css", "./style.css")
	router.Static("/picture", "./picture")

	// 加載 HTML 模板
	router.LoadHTMLGlob("./webpage/**/*")
	// router.GET("/components/header.html", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "header.html", nil)
	// })

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)
	router.POST("/SendLink", service.SendLink)

	// 2024.10.03 傳送連結到user mail後，驗證token並顯示更改密碼頁面
	router.GET("/webpage/login/reset_password", service.ResetPasswordPage)
	router.POST("/reset-password", service.ResetPassword)

	// 未登入的使用者
	// noUser := models.User{
	// 	Id:            "未登入",
	// 	Password_hash: "未登入",
	// 	Name:          "未登入",
	// 	Email:         "未登入",
	// 	Permission:    "1",
	// }

	// 針對學生的路由
	router.GET("/webpage/Student/student.html", service.RoleMiddleware("1"))
	router.GET("/Announcements", service.RenderAnnouncement)

	router.GET("/Attendance_record", auth.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "Attendance_record.html", nil)
	})

	// 針對管理員的路由
	router.GET("/webpage/manager/manager.html", service.RoleMiddleware("2"))

	// 針對教授的路由
	router.GET("/webpage/Professer/professer.html", service.RoleMiddleware("3"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "professer.html", nil)
	})

	// 測試部分
	router.GET("/test", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "Attendance_record.html", nil)
	})

	router.GET("/manager", func(c *gin.Context) {
		c.HTML(200, "manager.html", nil)
	})

	// 將通配符路由放在最後
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil) // 可以設置404頁面
	})

	router.Run(":8080")
}
