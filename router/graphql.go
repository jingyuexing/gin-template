package router

import (
	"template/api/gql"

	"github.com/gin-gonic/gin"
)

func graphQLRouterInit(public *gin.RouterGroup, private *gin.RouterGroup) {
	publicAPI := public.Group("v1")

	graphql := publicAPI.Group("graphql")
	graphql.POST("",gql.GraphQLAPI)
	graphql.GET("",gql.GraphQLAPI)

}
