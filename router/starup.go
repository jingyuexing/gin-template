package router

import (
	"fmt"
	api "template/api/v1"
	"template/global"
	"template/middleware"

	"github.com/gin-gonic/gin"
)
var API = api.AppAPI
func Routers() *gin.Engine {
	gin.SetMode(global.Config.Env.GinMode)

	fmt.Println("current Mode is: (" + gin.Mode() + ")")

	RootRouter := gin.Default()

	RootRouter.Use(
		middleware.LimitRate(
			middleware.WithQPS(1),
			middleware.WithBurst(32),
		),
		middleware.Authrize(
			middleware.WithSkipPaths(
				"/",
				"/api",
			),
		),
	)

	publicRouter := RootRouter.Group("/api")
	//
	privateRouter := RootRouter.Group("/api")

	var routerInit RouterBootList = RouterBootList{
		// your other module router in here

		authRouterInit,
		userRouterInit,
		graphQLRouterInit,
	}

	for _, RB := range routerInit {
		RB(publicRouter, privateRouter)
	}
	return RootRouter
}
