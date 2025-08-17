package user

import (
	"fmt"
	"strings"
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

func (d *UserDao) Search(values map[string]string,page,size int) ([]*model.UserModel, error) {
	result := []*model.UserModel{}
	statement := db.Model(&model.UserModel{})
	inited := false

	// 默认分页参数
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10 // 设置默认页大小为10
	}

	for key, value := range values {
		// 排除分页相关的字段
		if key == "page" || key == "size" {
			continue
		}

		// 基本条件：查找匹配的字段
		condition := key + " = ?"
		// 处理 "!" 表示不等于
		if strings.Contains(fmt.Sprintf("%v", value), "!") {
			condition = key + " != ?"
			value = strings.Replace(fmt.Sprintf("%v", value), "!", "", 1)
		}

		// 处理 LIKE 操作
		if strings.Contains(fmt.Sprintf("%v", value), "%") {
			// value like %"name"
			value = strings.Replace(fmt.Sprintf("%v", value), "%", "", 1)
			statement = statement.Where(key+" LIKE ?", "%"+value+"%")
			continue
		}

		// 第一次初始化条件
		if !inited {
			statement = statement.Where(condition, value)
			statement = statement.Offset((page - 1) * size).Limit(size)
			inited = true
			continue
		}

		// 处理 OR 条件
		statement = statement.Or(condition, value)
	}

	err := statement.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

