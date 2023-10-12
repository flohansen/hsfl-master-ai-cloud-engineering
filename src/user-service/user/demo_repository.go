package user

import (
	"errors"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
)

type DemoRepository struct {
	users map[uint64]*model.User
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{users: make(map[uint64]*model.User)}
}

func (repo *DemoRepository) Create(user *model.User) (*model.User, error) {
	_, found := repo.users[user.Id]
	if found {
		return nil, errors.New("user already exists")
	}
	repo.users[user.Id] = user

	return user, nil
}

func (repo *DemoRepository) Delete(user *model.User) error {
	_, found := repo.users[user.Id]
	if found {
		delete(repo.users, user.Id)
		return nil
	}

	return errors.New("user could not be deleted")
}

func (repo *DemoRepository) FindByEmail(email string) (*model.User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user could not be found")
}

func (repo *DemoRepository) FindById(id uint64) (*model.User, error) {
	user, found := repo.users[id]
	if found {
		return user, nil
	}

	return nil, errors.New("user could not be found")
}

func (repo *DemoRepository) Update(user *model.User) (*model.User, error) {
	user, foundError := repo.FindById(user.Id)

	if foundError != nil {
		return nil, errors.New("user can not be updated")
	}

	return user, nil
}
