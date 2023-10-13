package user

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"

type Repository interface {
	Create(*model.User) (*model.User, error)
	Delete(*model.User) error
	FindByEmail(email string) (*model.User, error)
	FindById(id uint64) (*model.User, error)
	Update(*model.User) (*model.User, error)
}
