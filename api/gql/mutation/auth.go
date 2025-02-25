package mutation

import (
	"template/api/gql/types"
	"template/core"
	"template/dto"
	"template/service"

	"github.com/graphql-go/graphql"
)


var LoginMutation = &graphql.Field{
	Type: types.AuthTokenType,
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		userLoginInput := &dto.UserLoginDTO{
			Email:    p.Args["email"].(string),
			Password: p.Args["password"].(string),
		}
		err := core.ValidateParams(userLoginInput)
		if err != nil {
			return nil, err
		}
		user,err := service.ServiceBoot.User.GetByGID("")
		if err != nil {
			return nil,err
		}
		access_token,refresh,err := service.ServiceBoot.Auth.Token(
			user.Name,
			user.GID,
			user.ID,
		)
		if err != nil {
			return nil,err
		}
		return map[string]any{
			"access_token":  access_token,
			"refresh_token": refresh,
		}, nil
	},
}
var RefreshMutation = &graphql.Field{
	Type: types.AuthTokenType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		refresh := ""
		return map[string]any{
			"refresh_token": refresh,
		}, nil
	},
}

