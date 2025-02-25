package gql

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var graphQLHandler *handler.Handler

func init() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery, // Query
		Mutation: rootMutation, // Mutation
	})
	if err != nil {
		panic("\nfailed to create schema: " + err.Error())
	}
	graphQLHandler = handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		Playground: false, // this is debug playground
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			return map[string]interface{}{
				"headers":r.Header,
			}
		},
	})
}

func GraphQLAPI(ctx *gin.Context) {
	graphQLHandler.ServeHTTP(ctx.Writer, ctx.Request)
}
