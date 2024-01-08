package userShoppingListEntry

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"

type Repository interface {
	Create(*model.UserShoppingListEntry) (*model.UserShoppingListEntry, error)
	FindByIds(shoppingListId uint64, productId uint64) (*model.UserShoppingListEntry, error)
	FindAll(shoppingListId uint64) ([]*model.UserShoppingListEntry, error)
	Update(list *model.UserShoppingListEntry) (*model.UserShoppingListEntry, error)
	Delete(*model.UserShoppingListEntry) error
}

const (
	ErrorEntryNotFound      = "entry could not be found"
	ErrorEntryUpdate        = "entry can not be updated"
	ErrorEntryDeletion      = "entry could not be deleted"
	ErrorEntryAlreadyExists = "entry already exists"
)
