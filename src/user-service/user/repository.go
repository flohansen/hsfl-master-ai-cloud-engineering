package user

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"

type Repository interface {
	Create(*model.User) (*model.User, error)
	Delete(*model.User) error
	FindAll() ([]*model.User, error)
	FindAllByRole(role *model.Role) ([]*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindById(id uint64) (*model.User, error)
	Update(*model.User) (*model.User, error)
}

const (
	ErrorUserList          = "user list not available"
	ErrorUserNotFound      = "user could not be found"
	ErrorUserUpdate        = "user can not be updated"
	ErrorUserDeletion      = "user could not be deleted"
	ErrorUserAlreadyExists = "user already exists"
)
