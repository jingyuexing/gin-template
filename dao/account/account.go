package account

import (
	"template/global"
	"template/model"
)

var db = global.DB

type IAccount interface {
    Create(account *model.AccountModel) error
    Delete(gid string) error
    Update(account *model.AccountModel) error
    FindByGID(gid string) (*model.AccountModel, error)
    List(page, size int) ([]*model.AccountModel, error)
}

type AccountDao struct {}

func (d *AccountDao) Create(account *model.AccountModel) error {
    return db.Create(account).Error
}

func (d *AccountDao) Delete(gid string) error {
    return db.Where("gid = ?", gid).Delete(&model.AccountModel{}).Error
}

func (d *AccountDao) Search(values map[string]any) (*model.AccountModel, error) {
	var account model.AccountModel
	err := db.Where(values).Find(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (d *AccountDao) Update(account *model.AccountModel) error {
    return db.Where("gid = ?", account.GID).Updates(account).Error
}

func (d *AccountDao) FindByGID(gid string) (*model.AccountModel, error) {
    var account model.AccountModel
    err := db.Where("gid = ?", gid).First(&account).Error
    return &account, err
}

func (d *AccountDao) List(page, size int) ([]*model.AccountModel, error) {
    var accounts []*model.AccountModel
    err := db.Offset((page - 1) * size).Limit(size).Find(&accounts).Error
    return accounts, err
}


