package router

import "github.com/gin-gonic/gin"

func userRouterInit(public *gin.RouterGroup, private *gin.RouterGroup){
	publicRouter := public.Group("v1")
	privateRouter := private.Group("v1")
	publicUser := publicRouter.Group("user")
	privateUser := privateRouter.Group("user")

	userController := API.User
	publicUser.POST("create", userController.CreateUser)

	privateUser.GET("list", userController.List)
	privateUser.PUT("update", userController.Update)
	privateUser.DELETE("delete", userController.Delete)
}
