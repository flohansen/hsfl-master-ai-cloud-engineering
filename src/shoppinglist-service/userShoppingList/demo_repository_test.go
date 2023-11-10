package userShoppingList

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
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

	list := &model.UserShoppingList{
		Id:        1,
		UserId:    1,
		Completed: false,
	}

	t.Run("Create shopping list with success", func(t *testing.T) {
		createdList, err := repo.Create(list)
		if err != nil {
			t.Error(err)
		}

		if createdList.Id != list.Id {
			t.Errorf("Expected created shopping list to have ID %d, but got %d", list.Id, createdList.Id)
		}
	})

	t.Run("Attempt to create a duplicate shopping list", func(t *testing.T) {
		_, err := repo.Create(list)
		if err == nil {
			t.Error("Expected an error for duplicate shopping list creation")
		}
	})
}

func TestDemoRepository_Delete(t *testing.T) {
	repo := NewDemoRepository()

	list := &model.UserShoppingList{
		Id:        1,
		UserId:    1,
		Completed: false,
	}

	_, _ = repo.Create(list)

	t.Run("Delete an existing shopping list", func(t *testing.T) {
		err := repo.Delete(list)
		if err != nil {
			t.Error(err)
		}

		// Ensure the list is deleted
		_, err = repo.findById(list.Id)
		if err == nil {
			t.Error("Expected the shopping list to be deleted")
		}
	})

	t.Run("Attempt to delete a non-existing shopping list", func(t *testing.T) {
		fakeList := &model.UserShoppingList{
			Id:        42,
			UserId:    2,
			Completed: false,
		}

		err := repo.Delete(fakeList)
		if err == nil {
			t.Error("Expected an error for deleting a non-existing shopping list")
		}
	})
}

func TestDemoRepository_Update(t *testing.T) {
	repo := NewDemoRepository()

	list := &model.UserShoppingList{
		Id:          1,
		UserId:      1,
		Description: "Update list description",
		Completed:   true,
	}

	_, _ = repo.Create(list)

	t.Run("Update an existing shopping list", func(t *testing.T) {
		updatedList, err := repo.Update(list)
		if err != nil {
			t.Error(err)
		}

		if updatedList.Description != list.Description {
			t.Errorf("Expected updated shopping list to have description %s, but got %s",
				list.Description, updatedList.Description)
		}

		if updatedList.Completed != list.Completed {
			t.Errorf("Expected updated shopping list to have completed value %t, but got %t",
				list.Completed, updatedList.Completed)
		}
	})

	t.Run("Attempt to update a non-existing shopping list", func(t *testing.T) {
		fakeList := &model.UserShoppingList{
			Id:        42,
			UserId:    2,
			Completed: true,
		}

		_, err := repo.Update(fakeList)
		if err == nil {
			t.Error("Expected an error for updating a non-existing shopping list")
		}
	})
}

func TestDemoRepository_FindAllById(t *testing.T) {
	repo := NewDemoRepository()

	lists := []*model.UserShoppingList{
		{
			Id:        1,
			UserId:    1,
			Completed: false,
		},
		{
			Id:        2,
			UserId:    2,
			Completed: false,
		},
		{
			Id:        3,
			UserId:    1,
			Completed: true,
		},
	}

	for _, list := range lists {
		_, _ = repo.Create(list)
	}

	t.Run("Find all lists for a user with existing lists", func(t *testing.T) {
		userLists, err := repo.FindAllById(1)
		if err != nil {
			t.Error(err)
		}

		if len(userLists) != 2 {
			t.Errorf("Expected 2 lists for user 1, but got %d", len(userLists))
		}
	})

	t.Run("Find all lists for a user with no existing lists", func(t *testing.T) {
		_, err := repo.FindAllById(42)
		if err == nil {
			t.Error("Expected an error for user with no lists")
		}
	})
}

func TestDemoRepository_FindByIds(t *testing.T) {
	repo := NewDemoRepository()

	lists := []*model.UserShoppingList{
		{
			Id:        1,
			UserId:    1,
			Completed: false,
		},
		{
			Id:        2,
			UserId:    2,
			Completed: false,
		},
	}

	for _, list := range lists {
		_, _ = repo.Create(list)
	}

	t.Run("Find list by IDs with existing list", func(t *testing.T) {
		foundList, err := repo.FindByIds(1, 1)
		if err != nil {
			t.Error(err)
		}

		if foundList.Id != 1 {
			t.Errorf("Expected list with ID 1, but got list with ID %d", foundList.Id)
		}
	})

	t.Run("Find list by IDs with non-existing list", func(t *testing.T) {
		_, err := repo.FindByIds(1, 42)
		if err == nil {
			t.Error("Expected an error for non-existing list")
		}
	})
}

func TestDemoRepository_findNextAvailableID(t *testing.T) {
	repo := NewDemoRepository()

	lists := []*model.UserShoppingList{
		{
			Id:        1,
			UserId:    1,
			Completed: false,
		},
		{
			Id:        2,
			UserId:    2,
			Completed: false,
		},
	}

	for _, list := range lists {
		_, _ = repo.Create(list)
	}

	t.Run("Next available ID when there are no gaps", func(t *testing.T) {
		nextID := repo.findNextAvailableID()
		if nextID != 3 {
			t.Errorf("Expected next available ID to be 4, but got %d", nextID)
		}
	})

	t.Run("Next available ID with gaps in IDs", func(t *testing.T) {
		delete(repo.shoppingLists, 2)
		nextID := repo.findNextAvailableID()
		if nextID != 2 {
			t.Errorf("Expected next available ID to be 2, but got %d", nextID)
		}
	})
}
