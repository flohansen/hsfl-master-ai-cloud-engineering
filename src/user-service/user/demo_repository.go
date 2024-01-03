package user

import (
	"errors"
	"sort"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
)

type DemoRepository struct {
	users map[uint64]*model.User
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{users: make(map[uint64]*model.User)}
}

func (repo *DemoRepository) Create(user *model.User) (*model.User, error) {
	var userId uint64
	if user.Id == 0 {
		userId = repo.findNextAvailableID()
	} else {
		userId = user.Id
	}

	_, found := repo.users[userId]
	if found {
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	repo.users[userId] = user

	return user, nil
}

func (repo *DemoRepository) Delete(user *model.User) error {
	_, found := repo.users[user.Id]
	if found {
		delete(repo.users, user.Id)
		return nil
	}

	return errors.New(ErrorUserDeletion)
}

func (repo *DemoRepository) FindAll() ([]*model.User, error) {
	if repo.users != nil {
		r := make([]*model.User, 0, len(repo.users))
		for _, v := range repo.users {
			r = append(r, v)
		}

		sort.Slice(r, func(i, j int) bool {
			return r[i].Name < r[j].Name
		})
		return r, nil
	}

	return nil, errors.New(ErrorUserList)
}

func (repo *DemoRepository) FindAllByRole(role *model.Role) ([]*model.User, error) {
	if repo.users != nil {
		r := make([]*model.User, 0, len(repo.users))
		for _, user := range repo.users {
			if user.Role == *role {
				r = append(r, user)
			}
		}

		sort.Slice(r, func(i, j int) bool {
			return r[i].Name < r[j].Name
		})
		return r, nil
	}

	return nil, errors.New(ErrorUserList)
}

func (repo *DemoRepository) FindByEmail(email string) (*model.User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New(ErrorUserNotFound)
}

func (repo *DemoRepository) FindById(id uint64) (*model.User, error) {
	user, found := repo.users[id]
	if found {
		return user, nil
	}

	return nil, errors.New(ErrorUserNotFound)
}

func (repo *DemoRepository) Update(user *model.User) (*model.User, error) {
	existingUser, foundError := repo.FindById(user.Id)

	if foundError != nil {
		return nil, errors.New(ErrorUserUpdate)
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Role = user.Role

	return existingUser, nil
}

func (repo *DemoRepository) findNextAvailableID() uint64 {
	var maxID uint64
	for id := range repo.users {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}
