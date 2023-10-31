package userShoppingList

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/userShoppingList/model"

type Repository interface {
	Create(*model.UserShoppingList) (*model.UserShoppingList, error)
	Delete(*model.UserShoppingList) error
	FindById(userId uint64) (*model.UserShoppingList, error)
}

const (
	ErrorListNotFound = "list could not be found"
)
