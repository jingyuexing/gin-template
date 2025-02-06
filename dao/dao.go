package dao

import (
	"template/dao/account"
	"template/dao/user"
)


type Dao struct {
	User user.UserDao
	Account account.AccountDao
}

var APIDao = new(Dao)