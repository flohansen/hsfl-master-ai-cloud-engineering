package user

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"

type Repository interface {
	Migrate() error
	Create([]*model.DbUser) error
	Update(username string, user *model.UpdateUser) error
	FindAll() ([]*model.DbUser, error)
	FindByEmail(email string) ([]*model.DbUser, error)
	FindByUsername(username string) ([]*model.DbUser, error)
	Delete([]*model.DbUser) error
}
