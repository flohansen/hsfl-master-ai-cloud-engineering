//go:build integration
// +build integration

package userShoppingList

import (
	"context"
	_ "github.com/rqlite/gorqlite/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"reflect"
	"testing"
	"time"
)

const TestPort = "7002"

func TestIntegrationRQLiteRepository(t *testing.T) {
	container, err := prepareIntegrationTestRQLiteDatabase()
	if err != nil {
		t.Error(err)
	}

	rqliteRepository := NewRQLiteRepository("http://localhost:" + TestPort + "/?disableClusterDiscovery=true")

	t.Run("TestIntegrationRQLiteRepository_Create", func(t *testing.T) {
		list := model.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "New Shoppinglist",
			Completed:   false,
		}

		t.Run("Create list with success", func(t *testing.T) {
			_, err = rqliteRepository.Create(&list)
			if err != nil {
				t.Error(err)
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindAllById", func(t *testing.T) {
		userId := uint64(2)

		lists := []*model.UserShoppingList{
			{
				Id:        1,
				UserId:    userId,
				Completed: false,
			},
			{
				Id:        2,
				UserId:    userId,
				Completed: false,
			},
			{
				Id:        3,
				UserId:    userId,
				Completed: true,
			},
		}

		for _, list := range lists {
			rqliteRepository.Create(list)
		}

		t.Run("Successfully fetch all lists from user", func(t *testing.T) {
			fetchedLists, err := rqliteRepository.FindAllById(userId)

			if err != nil {
				t.Error(err)
			}

			if len(fetchedLists) != len(lists) {
				t.Errorf("Unexpected list count. Expected %d, got %d", len(lists), len(fetchedLists))
			}

			if !reflect.DeepEqual(lists, fetchedLists) {
				t.Error("Fetched lists do not match expected lists")
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindById", func(t *testing.T) {
		list := model.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "New List",
			Completed:   false,
		}

		rqliteRepository.Create(&list)

		t.Run("Successfully fetch list", func(t *testing.T) {
			fetchedList, err := rqliteRepository.FindById(list.Id)
			if err != nil {
				t.Errorf("Can't find expected list with id %d: %v", list.Id, err)
			}

			if !reflect.DeepEqual(list, *fetchedList) {
				t.Error("Fetched list does not match original list")
			}
		})

		t.Run("Fail to fetch list (list not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindById(10)
			if err == nil {
				t.Errorf("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindByIds", func(t *testing.T) {
		list := model.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "New List",
			Completed:   false,
		}

		rqliteRepository.Create(&list)

		t.Run("Successfully fetch list", func(t *testing.T) {
			fetchedList, err := rqliteRepository.FindByIds(list.Id, list.UserId)
			if err != nil {
				t.Errorf("Can't find expected list with id %d: %v", list.Id, err)
			}

			if !reflect.DeepEqual(list, *fetchedList) {
				t.Error("Fetched list does not match original list")
			}
		})

		t.Run("Fail to fetch list (list not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindByIds(10, 20)
			if err == nil {
				t.Errorf("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_Update", func(t *testing.T) {
		list := model.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Old list description",
			Completed:   true,
		}

		rqliteRepository.Create(&list)

		t.Run("Update list with success", func(t *testing.T) {
			changedList := model.UserShoppingList{
				Id:          1,
				UserId:      1,
				Description: "Update list description",
				Completed:   true,
			}

			updatedList, err := rqliteRepository.Update(&changedList)
			if reflect.DeepEqual(changedList, updatedList) || err != nil {
				t.Error("Failed to list list")
			}
		})

		t.Run("Update list with fail (list not found)", func(t *testing.T) {
			changedList := model.UserShoppingList{
				Id:          10,
				UserId:      2,
				Description: "Update list description",
				Completed:   false,
			}

			_, err := rqliteRepository.Update(&changedList)
			if err == nil {
				t.Error("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_Delete", func(t *testing.T) {
		listToDelete := model.UserShoppingList{
			Id:        1,
			UserId:    1,
			Completed: false,
		}

		rqliteRepository.Create(&listToDelete)

		t.Run("Delete list with success", func(t *testing.T) {
			err = rqliteRepository.Delete(&listToDelete)
			if err != nil {
				t.Errorf("Failed to delete list with id %d", listToDelete.Id)
			}
		})

		t.Run("Delete list with fail (list not found)", func(t *testing.T) {
			notExistingList := model.UserShoppingList{
				Id:        10,
				UserId:    1,
				Completed: false,
			}
			err := rqliteRepository.Delete(&notExistingList)
			if err == nil {
				t.Error("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Cleanup(func() {
		err = container.Stop(context.Background(), nil)
		if err != nil {
			return
		}
	})
}

func prepareIntegrationTestRQLiteDatabase() (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        "rqlite/rqlite:8.15.0",
		ExposedPorts: []string{TestPort + ":4001/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("4001/tcp"),
			wait.ForLog(`.*HTTP API available at.*`).AsRegexp(),
			wait.ForLog(".*entering leader state.*").AsRegexp()).
			WithStartupTimeoutDefault(120 * time.Second).
			WithDeadline(360 * time.Second),
	}

	return testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
}
