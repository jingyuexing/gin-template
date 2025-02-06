package api

import (
	"template/api/auth"
	"template/api/user"
)

type API struct {
	Auth auth.AuthController
	User user.UserController
}

var AppAPI = new(API)