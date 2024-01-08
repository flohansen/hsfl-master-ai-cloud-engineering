package userShoppingListEntry

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"testing"
)

func TestNewDemoRepository(t *testing.T) {
	t.Run("Demo repository correctly initialized", func(t *testing.T) {
		repo := NewDemoRepository()
		if repo == nil {
			t.Error("NewDemoRepository returned nil")
		}
	})
}

func TestDemoRepository_Create(t *testing.T) {
	repo := NewDemoRepository()

	entry := &model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          2,
		Note:           "This is very important",
		Checked:        false,
	}

	t.Run("Create entry with success", func(t *testing.T) {
		createdEntry, err := repo.Create(entry)
		if err != nil {
			t.Error(err)
		}

		if createdEntry.ShoppingListId != entry.ShoppingListId || createdEntry.ProductId != entry.ProductId {
			t.Errorf("Expected created entry to have ShoppingListId %d and ProductId %d, but got ShoppingListId %d and ProductId %d",
				entry.ShoppingListId, entry.ProductId, createdEntry.ShoppingListId, createdEntry.ProductId)
		}
	})

	t.Run("Attempt to create a duplicate entry", func(t *testing.T) {
		_, err := repo.Create(entry)
		if err == nil {
			t.Error("Expected an error for duplicate entry creation")
		}
	})
}

func TestDemoRepository_FindAll(t *testing.T) {
	repo := NewDemoRepository()

	entries := []*model.UserShoppingListEntry{
		{
			ShoppingListId: 1,
			ProductId:      1,
			Count:          2,
			Note:           "This is very important",
			Checked:        false,
		},
		{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          5,
			Note:           "I really want this",
			Checked:        false,
		},
		{
			ShoppingListId: 2,
			ProductId:      1,
			Count:          9,
			Note:           "Please get this",
			Checked:        false,
		},
	}

	for _, entry := range entries {
		_, _ = repo.Create(entry)
	}

	t.Run("Find all entries for a shopping list with existing entries", func(t *testing.T) {
		shoppingListId := uint64(1)
		foundEntries, err := repo.FindAll(shoppingListId)
		if err != nil {
			t.Error(err)
		}

		if len(foundEntries) != 2 {
			t.Errorf("Expected 2 entries for shopping list %d, but got %d", shoppingListId, len(foundEntries))
		}
	})

	t.Run("Find all entries for a shopping list with no existing entries", func(t *testing.T) {
		shoppingListId := uint64(3)
		_, err := repo.FindAll(shoppingListId)
		if err == nil {
			t.Error("Expected an error for shopping list with no entries")
		}
	})
}

func TestDemoRepository_FindByIds(t *testing.T) {
	repo := NewDemoRepository()

	entry := &model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          2,
		Note:           "This is very important",
		Checked:        false,
	}

	_, _ = repo.Create(entry)

	t.Run("Find entry by IDs with existing entry", func(t *testing.T) {
		foundEntry, err := repo.FindByIds(entry.ShoppingListId, entry.ProductId)
		if err != nil {
			t.Error(err)
		}

		if foundEntry.ShoppingListId != entry.ShoppingListId || foundEntry.ProductId != entry.ProductId {
			t.Errorf("Expected entry with ShoppingListId %d and ProductId %d, but got ShoppingListId %d and ProductId %d",
				entry.ShoppingListId, entry.ProductId, foundEntry.ShoppingListId, foundEntry.ProductId)
		}
	})

	t.Run("Find entry by IDs with non-existing entry", func(t *testing.T) {
		nonExistentEntry, err := repo.FindByIds(2, 2)
		if err == nil {
			t.Error("Expected an error for non-existing entry")
		}

		if nonExistentEntry != nil {
			t.Errorf("Expected non-existing entry to be nil, but got an entry")
		}
	})
}

func TestDemoRepository_Update(t *testing.T) {
	repo := NewDemoRepository()

	entry := &model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          3,
		Note:           "I really want this",
		Checked:        true,
	}

	_, _ = repo.Create(entry)

	t.Run("Update an existing entry", func(t *testing.T) {
		updatedEntry, err := repo.Update(entry)
		if err != nil {
			t.Error(err)
		}

		if updatedEntry.Count != entry.Count {
			t.Errorf("Expected updated shopping list entry to have count %d, but got count %d",
				entry.Count, updatedEntry.Count)
		}

		if updatedEntry.Note != entry.Note {
			t.Errorf("Expected updated shopping list entry to have note %s, but got note %s",
				entry.Note, updatedEntry.Note)
		}

		if updatedEntry.Checked != entry.Checked {
			t.Errorf("Expected updated shopping list entry to have checked value %t, but got note %t",
				entry.Checked, updatedEntry.Checked)
		}
	})

	t.Run("Attempt to update a non-existing entry", func(t *testing.T) {
		fakeEntry := &model.UserShoppingListEntry{
			ShoppingListId: 2,
			ProductId:      2,
			Count:          2,
			Note:           "This is very important",
			Checked:        false,
		}

		_, err := repo.Update(fakeEntry)
		if err == nil {
			t.Error("Expected an error for updating a non-existing entry")
		}
	})
}

func TestDemoRepository_Delete(t *testing.T) {
	repo := NewDemoRepository()

	entry := &model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          1,
		Note:           "This is very important to me",
		Checked:        false,
	}

	_, _ = repo.Create(entry)

	t.Run("Delete an existing entry", func(t *testing.T) {
		err := repo.Delete(entry)
		if err != nil {
			t.Error(err)
		}

		// Ensure the entry is deleted
		_, err = repo.FindByIds(entry.ShoppingListId, entry.ProductId)
		if err == nil {
			t.Error("Expected the entry to be deleted")
		}
	})

	t.Run("Attempt to delete a non-existing entry", func(t *testing.T) {
		fakeEntry := &model.UserShoppingListEntry{
			ShoppingListId: 2,
			ProductId:      2,
			Count:          2,
			Note:           "This is very important",
			Checked:        false,
		}

		err := repo.Delete(fakeEntry)
		if err == nil {
			t.Error("Expected an error for deleting a non-existing entry")
		}
	})
}
