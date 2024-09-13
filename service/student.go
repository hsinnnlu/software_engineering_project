package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StudentPage(c *gin.Context) {
	c.HTML(http.StatusOK, "student.html", nil)
}
