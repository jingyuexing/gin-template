package types

import "github.com/graphql-go/graphql"

var OauthType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Oauth",
	Fields: graphql.Fields{
		"client_id":&graphql.Field{Type: graphql.String},
		"client_secret": &graphql.Field{Type: graphql.String},
		"redirect_uri": &graphql.Field{Type: graphql.String},
		"grant_type": &graphql.Field{Type: graphql.String},
		"scope":&graphql.Field{Type: graphql.String},
	},
})
var AuthTokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthToken",
	Fields: graphql.Fields{
		"access_token": &graphql.Field{
			Type: graphql.String,
		},
		"refresh_token": &graphql.Field{
			Type: graphql.String,
		},
	},
})
