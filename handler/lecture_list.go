package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/service"
)

func Lecturelisthandler(c *gin.Context){
	lectures, err := service.GetActiveLectures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lecture error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lecture": lectures})
}