package router

import "github.com/gin-gonic/gin"

func userRouterInit(public *gin.RouterGroup, private *gin.RouterGroup){
	publicRouter := public.Group("v1")
	privateRouter := private.Group("v1")
	publicUser := publicRouter.Group("user")
	privateUser := privateRouter.Group("user")

	// userController := API.User

	// publicUser.GET("list", userController.List)
	// publicUser.POST("create", userController.Create)
	// publicUser.PUT("update", userController.Update)
	// publicUser.DELETE("delete", userController.Delete)
}
