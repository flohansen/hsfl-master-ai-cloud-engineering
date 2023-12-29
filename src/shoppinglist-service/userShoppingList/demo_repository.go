package userShoppingList

import (
	"errors"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"sort"
)

type DemoRepository struct {
	shoppingLists map[uint64]*model.UserShoppingList
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{shoppingLists: make(map[uint64]*model.UserShoppingList)}
}

func (repo *DemoRepository) Create(shoppingList *model.UserShoppingList) (*model.UserShoppingList, error) {
	var listId uint64
	if shoppingList.Id == 0 {
		listId = repo.findNextAvailableID()
		shoppingList.Id = listId
	} else {
		listId = shoppingList.Id
	}

	_, found := repo.shoppingLists[listId]
	if found {
		return nil, errors.New(ErrorListAlreadyExists)
	}
	repo.shoppingLists[listId] = shoppingList

	return shoppingList, nil
}

func (repo *DemoRepository) Delete(shoppingList *model.UserShoppingList) error {
	_, found := repo.shoppingLists[shoppingList.Id]
	if found {
		delete(repo.shoppingLists, shoppingList.Id)
		return nil
	}

	return errors.New(ErrorListDeletion)
}

func (repo *DemoRepository) Update(shoppingList *model.UserShoppingList) (*model.UserShoppingList, error) {
	existingShoppingList, foundError := repo.findById(shoppingList.Id)

	if foundError != nil {
		return nil, errors.New(ErrorListUpdate)
	}

	existingShoppingList.Description = shoppingList.Description
	existingShoppingList.Completed = shoppingList.Completed

	return existingShoppingList, nil
}

func (repo *DemoRepository) FindAllById(userId uint64) ([]*model.UserShoppingList, error) {
	lists := []*model.UserShoppingList{}

	for _, shoppingList := range repo.shoppingLists {
		if shoppingList.UserId == userId {
			lists = append(lists, shoppingList)
		}
	}

	if len(lists) == 0 {
		return nil, errors.New(ErrorListNotFound)
	}

	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Description < lists[j].Description
	})

	return lists, nil
}

func (repo *DemoRepository) findById(Id uint64) (*model.UserShoppingList, error) {
	shoppingList, found := repo.shoppingLists[Id]
	if found {
		return shoppingList, nil
	}

	return nil, errors.New(ErrorListNotFound)
}

func (repo *DemoRepository) FindByIds(userId uint64, listId uint64) (*model.UserShoppingList, error) {
	for _, shoppingList := range repo.shoppingLists {
		if shoppingList.UserId == userId && shoppingList.Id == listId {
			return shoppingList, nil
		}
	}

	return nil, errors.New(ErrorListNotFound)
}

func (repo *DemoRepository) findNextAvailableID() uint64 {
	var maxID uint64
	for id := range repo.shoppingLists {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}
