// 目前沒有任何程式 call 到這個檔案，所以可以先不用管他

package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StudentPage(c *gin.Context) {
	c.HTML(http.StatusOK, "student.html", nil)
}
