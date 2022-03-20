package response

import (
	"net/http"

	"github.com/HideInBush7/go_im/status"
	"github.com/gin-gonic/gin"
)

func Abort(c *gin.Context, code int32) {
	AbortWithAll(c, code, status.Message(code), nil)
}

func AbortError(c *gin.Context, err error) {
	AbortWithAll(c, status.ERROR, err.Error(), nil)
}

func AbortSuccess(c *gin.Context, data interface{}) {
	AbortWithAll(c, status.SUCCESS, status.Message(status.SUCCESS), data)
}

func AbortWithAll(c *gin.Context, code int32, msg string, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	})
}
