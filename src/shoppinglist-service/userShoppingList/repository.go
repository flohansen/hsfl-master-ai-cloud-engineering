package userShoppingList

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"

type Repository interface {
	Create(*model.UserShoppingList) (*model.UserShoppingList, error)
	FindAllById(userId uint64) ([]*model.UserShoppingList, error)
	FindById(listId uint64) (*model.UserShoppingList, error)
	FindByIds(userId uint64, listId uint64) (*model.UserShoppingList, error)
	Update(list *model.UserShoppingList) (*model.UserShoppingList, error)
	Delete(*model.UserShoppingList) error
}

const (
	ErrorListNotFound      = "list could not be found"
	ErrorListUpdate        = "list can not be updated"
	ErrorListDeletion      = "list could not be deleted"
	ErrorListAlreadyExists = "list already exists"
)
