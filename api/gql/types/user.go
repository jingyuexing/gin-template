package types

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "user",
	Fields: graphql.Fields{
	},
})
