package user

import (
	"template/global"
	"template/model"
)

var db = global.DB

type IUser interface {
    Create(user *model.UserModel) error
    Delete(gid string) error
    Update(user *model.UserModel) error
    FindByGID(gid string) (*model.UserModel, error)
    List(page, size int) ([]*model.UserModel, error)
}

type UserDao struct {}

func (d *UserDao) Create(user *model.UserModel) error {
    return db.Create(user).Error
}

func (d *UserDao) Delete(gid string) error {
    return db.Where("gid = ?", gid).Delete(&model.UserModel{}).Error
}

func (d *UserDao) Update(user *model.UserModel) error {
    return db.Where("gid = ?", user.GID).Updates(user).Error
}

func (d *UserDao) FindByGID(gid string) (*model.UserModel, error) {
    var user model.UserModel
    err := db.Where("gid = ?", gid).First(&user).Error
    return &user, err
}

func (d *UserDao) List(page, size int) ([]*model.UserModel, error) {
    var users []*model.UserModel
    err := db.Offset((page - 1) * size).Limit(size).Find(&users).Error
    return users, err
}

