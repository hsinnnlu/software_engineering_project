package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/software_engineering_project/db"

	"github.com/hsinnnlu/software_engineering_project/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db.InitDB("./database.db")

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions(("mysession"), store))

	router.Static("/style.css", "./style.css")
	router.Static("/picture", "./picture")
	router.LoadHTMLFiles(
		"./reset_password.html",
	)
	router.LoadHTMLGlob(("./webpage/**/*"))

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)

	router.GET("/student", service.StudentPage)

	// router.POST("/test", service.TestDB)

	router.Run(":8080")
}
