package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/software_engineering_project/db"
	. "github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/service"
)

func main() {

	DB, err := db.InitDB()
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}
	defer DB.Close()

	InitDB()

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions(("mysession"), store))

	router.Static("/style.css", "./style.css")
	router.Static("picture", "./picture")
	router.LoadHTMLFiles(
		"./login.html",
		"./reset_password.html",
	)
	router.LoadHTMLGlob(("webpage/**/*"))

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)

	router.POST("/test", TestDB)

	router.Run(":8080")
}
