package userShoppingList

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"

type Repository interface {
	Create(*model.UserShoppingList) (*model.UserShoppingList, error)
	Delete(*model.UserShoppingList) error
	Update(list *model.UserShoppingList) (*model.UserShoppingList, error)
	FindById(userId uint64) (*model.UserShoppingList, error)
	FindAll() ([]*model.UserShoppingList, error)
}

const (
	ErrorListNotFound      = "list could not be found"
	ErrorListAlreadyExists = "list already exists"
)
