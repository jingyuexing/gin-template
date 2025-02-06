package router

import "github.com/gin-gonic/gin"

func authRouterInit(public *gin.RouterGroup, private *gin.RouterGroup) {
	publicRouter := public.Group("v1")
	privateRouter := private.Group("v1")
	publicAuth := publicRouter.Group("auth")
	privateAuth := privateRouter.Group("auth")

	authController := API.Auth
	
	privateAuth.GET("refresh_token", authController.Refresh)

	publicAuth.POST("login", authController.Login)
}