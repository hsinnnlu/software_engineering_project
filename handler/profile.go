package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/software_engineering_project/service"
)

func Announcehandler(c *gin.Context) {
	announces, err := service.GetAnnouncementList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Announces error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"announces": announces})
}
