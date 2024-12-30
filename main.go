package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/handler"
)

func main() {
	if err := db.InitDB("./database.db"); err != nil {
		panic(err)
	}

	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允許的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 是否允許跨域攜帶 Cookie
		MaxAge:           12 * time.Hour, // OPTIONS 請求的緩存時間
	}))

	r.POST("/login", handler.LoginHandler) // finished (back, front)
	r.POST("/sendemail", handler.ResetpasswdSendlinkHandler)

	resetpasswdtoken := r.Group("/")
	resetpasswdtoken.Use(handler.ResetpasswdVerifylinkHandler) // 尚未測試（之後再說）
	{
		resetpasswdtoken.POST("/resetpasswd", handler.ResetpasswdChangepasswd)
	}

	// 身份驗證相關路由
	authorized := r.Group("/")
	authorized.Use(handler.AuthMiddleware())
	{
		authorized.POST("/announce", handler.Announcehandler)
		authorized.POST("/lecture", handler.Lecturelisthandler)
		authorized.POST("/:lecture_id/:user_id/:sign-in")
		// authorized.GET("/userinfo", service.GetUserProfile)
		// authorized.GET("/userinfo", services.GetUserProfile)        // finished
		// authorized.POST("/addVocab", handlers.AddVocabularyHandler) // 新增單字 //finished
		// authorized.POST("/addFavorite", handlers.AddFavoriteVocab)
	}

	r.Run(":8888")
}
