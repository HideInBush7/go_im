package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {

}

// 允许跨域
func Cors(ctx *gin.Context) {
	method := ctx.Request.Method
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
	ctx.Set("content-type", "application/json")
	if method == "OPTIONS" {
		ctx.JSON(http.StatusOK, nil)
	}
	ctx.Next()
}
