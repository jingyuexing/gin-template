package router

import (
	"fmt"
	"template/api"
	"template/global"

	"github.com/gin-gonic/gin"
)
var API = api.AppAPI
func Routers() *gin.Engine {
	gin.SetMode(global.Config.Env.GinMode)

	fmt.Println("current Mode is: (" + gin.Mode() + ")")

	RootRouter := gin.Default()

	publicRouter := RootRouter.Group("/api")
	//
	privateRouter := RootRouter.Group("/api")

	var routerInit RouterBootList = RouterBootList{
		// your other module router in here
		authRouterInit,
		userRouterInit,
	}

	for _, RB := range routerInit {
		RB(publicRouter, privateRouter)
	}
	return RootRouter
}
