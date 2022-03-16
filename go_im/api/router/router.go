package router

import (
	"github.com/HideInBush7/go_im/api/handler"
	"github.com/HideInBush7/go_im/api/middleware"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors)

	userGroup := r.Group("/user")
	{
		userGroup.POST("/login")
		userGroup.POST("/register", handler.Register)
	}
	userGroup.Use(middleware.Auth)
	{
		userGroup.POST("/logout")
	}
	// sendGroup := r.Group("/")
	return r
}
