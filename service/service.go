package service

import (
	"template/service/account"
	"template/service/auth"
	"template/service/user"
)

type Service struct {
	User user.UserService
	Auth auth.AuthService
	Account account.AccountService
}

var ServiceBoot = new(Service)