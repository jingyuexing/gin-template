package api

import (
	"template/api/v1/auth"
	"template/api/v1/user"
)

type API struct {
	Auth auth.AuthController
	User user.UserController
}

var AppAPI = new(API)
