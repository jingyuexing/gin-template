package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	GID         string `json:"gid" gorm:"column:gid;index;comment:'全局唯一ID'"`
	Name        string `json:"name" gorm:"column:name;comment:'用户名'"`
	Email       string `json:"email" gorm:"column:email;comment:'邮箱地址'"`
	Description string `json:"description" gorm:"column:description;comment:'描述'"`
	Phone       string `json:"phone" gorm:"column:phone;comment:'手机号'"`
	Country     string `json:"country" gorm:"column:country_code;comment:'国家代码'"`
	Gender      string `json:"gender" gorm:"column:gender;comment:'性别'"`
	Avatar      string `json:"avatar" gorm:"column:avatar;comment:'头像'"`
	Follows     int64  `json:"follows" gorm:"column:follows;comment:'粉丝数'"`
	Following   int64  `json:"following" gorm:"column:following;comment:'关注数'"`
}

// User选项方法
func WithUserGID(gid string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.GID = gid
    }
}

func WithUserName(name string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Name = name
    }
}

func WithUserEmail(email string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Email = email
    }
}

func WithUserDescription(description string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Description = description
    }
}

func WithUserPhone(phone string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Phone = phone
    }
}

func WithUserCountry(country string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Country = country
    }
}

func WithUserGender(gender string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Gender = gender
    }
}

func WithUserAvatar(avatar string) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Avatar = avatar
    }
}

func WithUserFollows(follows int64) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Follows = follows
    }
}

func WithUserFollowing(following int64) ModelOption[UserModel] {
    return func(u *UserModel) {
        u.Following = following
    }
}
