package userShoppingList

import (
	"errors"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
)

type DemoRepository struct {
	shoppinglists map[uint64]*model.UserShoppingList
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{shoppinglists: make(map[uint64]*model.UserShoppingList)}
}

func (repo *DemoRepository) Create(shoppinglist *model.UserShoppingList) (*model.UserShoppingList, error) {
	var listId uint64
	if shoppinglist.Id == 0 {
		listId = repo.findNextAvailableID()
	} else {
		listId = shoppinglist.Id
	}

	_, found := repo.shoppinglists[listId]
	if found {
		return nil, errors.New(ErrorListAlreadyExists)
	}
	repo.shoppinglists[listId] = shoppinglist

	return shoppinglist, nil
}

func (repo *DemoRepository) Delete(shoppinglist *model.UserShoppingList) error {
	_, found := repo.shoppinglists[shoppinglist.Id]
	if found {
		delete(repo.shoppinglists, shoppinglist.Id)
		return nil
	}

	return errors.New(ErrorListDeletion)
}

func (repo *DemoRepository) Update(shoppinglist *model.UserShoppingList) (*model.UserShoppingList, error) {
	product, foundError := repo.findById(shoppinglist.Id)

	if foundError != nil {
		return nil, errors.New(ErrorListUpdate)
	}

	return product, nil
}

func (repo *DemoRepository) FindAllById(userId uint64) ([]*model.UserShoppingList, error) {
	lists := []*model.UserShoppingList{}

	for _, shoppingList := range repo.shoppinglists {
		if shoppingList.UserId == userId {
			lists = append(lists, shoppingList)
		}
	}

	if len(lists) == 0 {
		return nil, errors.New(ErrorListNotFound)
	}

	return lists, nil
}

func (repo *DemoRepository) findById(Id uint64) (*model.UserShoppingList, error) {
	product, found := repo.shoppinglists[Id]
	if found {
		return product, nil
	}

	return nil, errors.New(ErrorListNotFound)
}

func (repo *DemoRepository) FindByIds(userId uint64, listId uint64) (*model.UserShoppingList, error) {
	for _, shoppinglist := range repo.shoppinglists {
		if shoppinglist.UserId == userId && shoppinglist.Id == listId {
			println(shoppinglist)
			return shoppinglist, nil
		}
	}

	return nil, errors.New(ErrorListNotFound)
}

func (repo *DemoRepository) findNextAvailableID() uint64 {
	var maxID uint64
	for id := range repo.shoppinglists {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}
