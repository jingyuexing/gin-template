package model

import "gorm.io/gorm"

type AccountModel struct {
	gorm.Model
	GID        string `json:"gid" gorm:"column:gid;index;comment:'全局唯一id'"`
	Email string `json:"email" gorm:"column:email;index;comment:'邮箱地址'"`
	StrID      string `json:"stringID" gorm:"column:str_id;index"`
	Role       string `json:"role" gorm:"column:role"`
	Permission int    `json:"permission" gorm:"column:permission"`
	Password   string `json:"password" gorm:"column:password"`
}

func (AccountModel) TableName() string {
	return "account"
}

// Account选项方法
func WithAccountGID(gid string) ModelOption[AccountModel] {
    return func(a *AccountModel) {
        a.GID = gid
    }
}

func WithAccountStrID(strID string) ModelOption[AccountModel] {
    return func(a *AccountModel) {
        a.StrID = strID
    }
}

func WithAccountRole(role string) ModelOption[AccountModel] {
    return func(a *AccountModel) {
        a.Role = role
    }
}

func WithAccountPermission(permission int) ModelOption[AccountModel] {
    return func(a *AccountModel) {
        a.Permission = permission
    }
}

func WithAccountPassword(password string) ModelOption[AccountModel] {
    return func(a *AccountModel) {
        a.Password = password
    }
}