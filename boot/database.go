package boot

import (
	"template/global"
	"template/model"
)

func database_boot(){
	global.DB.AutoMigrate(
		&model.UserModel{},
		&model.AccountModel{},
	)
}