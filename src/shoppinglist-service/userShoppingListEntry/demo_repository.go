package userShoppingListEntry

import (
	"errors"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
)

type shoppingListEntryKey struct {
	ShoppingListID uint64
	ProductID      uint64
}

type DemoRepository struct {
	entries map[shoppingListEntryKey]*model.UserShoppingListEntry
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{entries: make(map[shoppingListEntryKey]*model.UserShoppingListEntry)}
}

func (repo *DemoRepository) Create(entry *model.UserShoppingListEntry) (*model.UserShoppingListEntry, error) {

	_, err := repo.FindByIds(entry.ShoppingListId, entry.ProductId)
	if err == nil {
		return nil, errors.New(ErrorEntryAlreadyExists)
	}
	key := shoppingListEntryKey{
		ShoppingListID: entry.ShoppingListId,
		ProductID:      entry.ProductId,
	}

	repo.entries[key] = entry

	return entry, nil
}

func (repo *DemoRepository) Delete(entry *model.UserShoppingListEntry) error {
	key := shoppingListEntryKey{
		ShoppingListID: entry.ShoppingListId,
		ProductID:      entry.ProductId,
	}

	_, exists := repo.entries[key]
	if !exists {
		return errors.New(ErrorEntryDeletion)
	}

	delete(repo.entries, key)
	return nil
}

func (repo *DemoRepository) Update(entry *model.UserShoppingListEntry) (*model.UserShoppingListEntry, error) {
	existingEntry, foundError := repo.FindByIds(entry.ShoppingListId, entry.ProductId)

	if foundError != nil {
		return nil, errors.New(ErrorEntryUpdate)
	}

	existingEntry.Count = entry.Count
	existingEntry.Note = entry.Note
	existingEntry.Checked = entry.Checked

	return existingEntry, nil
}

func (repo *DemoRepository) FindByIds(shoppingListId, productId uint64) (*model.UserShoppingListEntry, error) {
	key := shoppingListEntryKey{
		ShoppingListID: shoppingListId,
		ProductID:      productId,
	}

	entry, exists := repo.entries[key]

	if !exists {
		return nil, errors.New(ErrorEntryNotFound)
	}

	return entry, nil
}

func (repo *DemoRepository) FindAll(shoppingListId uint64) ([]*model.UserShoppingListEntry, error) {
	entries := []*model.UserShoppingListEntry{}

	for key, entry := range repo.entries {
		if key.ShoppingListID == shoppingListId {
			entries = append(entries, entry)
		}
	}

	if len(entries) == 0 {
		return nil, errors.New(ErrorEntryNotFound)
	}

	return entries, nil
}
