package user

import (
	"template/dao"
	"template/dto"
	"template/internal/builtin"
	"template/model"
)

type IUserService interface {
	Create(user *model.UserModel) error
	Delete(gid string) error
	Update(user *model.UserModel) error
	GetByGID(gid string) (*model.UserModel, error)
	List(page, size int) ([]*model.UserModel, error)
}

type UserService struct {
}
var userDao = dao.APIDao.User

func (s *UserService) Create(user *model.UserModel) error {
	// 检查用户名是否存在
	existingUser, err := userDao.FindByGID(user.GID)
	if err == nil && existingUser != nil {
		return builtin.ErrUserNameExists
	}

	if err = userDao.Create(user); err != nil {
		return builtin.ErrInternalServer
	}
	return nil
}

func (s *UserService) Delete(gid string) error {
	// 检查用户是否存在
	existingUser, err := userDao.FindByGID(gid)
	if err != nil || existingUser == nil {
		return builtin.ErrUserNotFound
	}

	if err = userDao.Delete(gid); err != nil {
		return builtin.ErrDBDeleteFailed
	}
	return nil
}

func (s *UserService) Update(user *model.UserModel) error {
	// 检查用户是否存在
	existingUser, err := userDao.FindByGID(user.GID)
	if err != nil || existingUser == nil {
		return builtin.ErrUserNotFound
	}

	if err = userDao.Update(user); err != nil {
		return builtin.ErrDBUpdateFailed
	}
	return nil
}

func (s *UserService) GetByGID(gid string) (*model.UserModel, error) {
	user, err := userDao.FindByGID(gid)
	if err != nil {
		return nil, builtin.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) List(pag dto.Pagination) ([]*model.UserModel, error) {

	users, err := userDao.List(pag.Page, pag.Size)
	if err != nil {
		return nil, builtin.ErrInternalServer
	}
	return users, nil
}
