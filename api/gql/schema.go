package gql

import (
	"template/api/gql/mutation"

	"github.com/graphql-go/graphql"
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{},
	},
)

var rootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"login":mutation.LoginMutation,
			"refresh":mutation.RefreshMutation,
		},
	},
)
