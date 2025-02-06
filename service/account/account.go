package account

import (
	"errors"
	"template/dao"
	"template/dto"
	"template/internal/builtin"
	"template/model"

	"gorm.io/gorm"
)

type IAccountService interface {
	Create(account *model.AccountModel) error
	Delete(gid string) error
	Update(account *model.AccountModel) error
	GetByGID(gid string) (*model.AccountModel, error)
	List(page dto.Pagination) ([]*model.AccountModel, error)
}

type AccountService struct{}

var accountDao = dao.APIDao.Account

func (s *AccountService) Create(account *model.AccountModel) error {
	// 检查账号是否已存在
	existingAccount, err := accountDao.FindByGID(account.GID)
	if err == nil && existingAccount != nil {
		return builtin.ErrUserNameExists
	}

	if err = accountDao.Create(account); err != nil {
		return builtin.ErrDBInsertFailed
	}
	return nil
}

func (s *AccountService) Delete(gid string) error {
	// 检查账号是否存在
	existingAccount, err := accountDao.FindByGID(gid)
	if err != nil || existingAccount == nil {
		return builtin.ErrUserNotFound
	}

	if err = accountDao.Delete(gid); err != nil {
		return builtin.ErrDBDeleteFailed
	}
	return nil
}

func (s *AccountService) Update(account *model.AccountModel) error {
	// 检查账号是否存在
	existingAccount, err := accountDao.FindByGID(account.GID)
	if err != nil || existingAccount == nil {
		return builtin.ErrUserNotFound
	}

	if err = accountDao.Update(account); err != nil {
		return builtin.ErrDBUpdateFailed
	}
	return nil
}

func (s *AccountService) GetByGID(gid string) (*model.AccountModel, error) {
	account, err := accountDao.FindByGID(gid)
	if err != nil {
		return nil, builtin.ErrUserNotFound
	}
	return account, nil
}

func (s *AccountService) GetByEmail(email string) (*model.AccountModel, error) {
	account ,err := accountDao.Search(map[string]any{
		"email":email,
	})
	if err != nil {
		if !errors.Is(err,gorm.ErrRecordNotFound){
			return nil, builtin.ErrDBQueryFailed
		}
		return nil,builtin.ErrUserNotFound
	}
	return account,nil
}

func (s *AccountService) List(pag dto.Pagination) ([]*model.AccountModel, error) {
	accounts, err := accountDao.List(pag.Page, pag.Size)
	if err != nil {
		return nil, builtin.ErrDBQueryFailed
	}
	return accounts, nil
}