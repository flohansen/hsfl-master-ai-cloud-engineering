//go:build integration
// +build integration

package userShoppingListEntry

import (
	"context"
	_ "github.com/rqlite/gorqlite/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"reflect"
	"testing"
	"time"
)

const TestPort = "7012"

func TestIntegrationRQLiteRepository(t *testing.T) {
	container, err := prepareIntegrationTestRQLiteDatabase()
	if err != nil {
		t.Error(err)
	}

	rqliteRepository := NewRQLiteRepository("http://localhost:" + TestPort + "/?disableClusterDiscovery=true")

	t.Run("TestIntegrationRQLiteRepository_Create", func(t *testing.T) {
		entry := model.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      1,
			Count:          2,
			Note:           "This is very important",
			Checked:        false,
		}

		t.Run("Create entry with success", func(t *testing.T) {
			_, err = rqliteRepository.Create(&entry)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("Can't create entries with duplicate ean", func(t *testing.T) {
			_, err = rqliteRepository.Create(&entry)
			if err.Error() != ErrorEntryAlreadyExists {
				t.Error(err)
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindAll", func(t *testing.T) {
		shoppingListId := uint64(1)

		entries := []*model.UserShoppingListEntry{
			{
				ShoppingListId: shoppingListId,
				ProductId:      1,
				Count:          2,
				Note:           "This is very important",
				Checked:        false,
			},
			{
				ShoppingListId: shoppingListId,
				ProductId:      2,
				Count:          5,
				Note:           "I really want this",
				Checked:        false,
			},
			{
				ShoppingListId: shoppingListId,
				ProductId:      3,
				Count:          9,
				Note:           "Please get this",
				Checked:        false,
			},
		}

		for _, entry := range entries {
			rqliteRepository.Create(entry)
		}

		t.Run("Successfully fetch all entries", func(t *testing.T) {
			fetchedEntries, err := rqliteRepository.FindAll(shoppingListId)

			if err != nil {
				t.Error("Can't fetch entries")
			}

			if len(fetchedEntries) != len(entries) {
				t.Errorf("Unexpected entry count. Expected %d, got %d", len(entries), len(fetchedEntries))
			}

			if !reflect.DeepEqual(entries, fetchedEntries) {
				t.Error("Fetched entries do not match expected entries")
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindByIds", func(t *testing.T) {
		entry := model.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      1,
			Count:          2,
			Note:           "This is very important",
			Checked:        false,
		}

		rqliteRepository.Create(&entry)

		t.Run("Successfully fetch entry", func(t *testing.T) {
			fetchedProduct, err := rqliteRepository.FindByIds(entry.ShoppingListId, entry.ProductId)
			if err != nil {
				t.Errorf("Can't find expected entry with shoppingListId %d and productId %d: %v", entry.ShoppingListId, entry.ProductId, err)
			}

			if !reflect.DeepEqual(&entry, fetchedProduct) {
				t.Error("Fetched entry does not match original entry")
			}
		})

		t.Run("Fail to fetch entry (entry not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindByIds(2, 2)
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
		entry := model.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      1,
			Count:          3,
			Note:           "I really want this",
			Checked:        true,
		}

		rqliteRepository.Create(&entry)

		t.Run("Update entry with success", func(t *testing.T) {
			changedEntry := model.UserShoppingListEntry{
				ShoppingListId: 1,
				ProductId:      1,
				Count:          3,
				Note:           "I dont want this anymore",
				Checked:        false,
			}

			updatedEntry, err := rqliteRepository.Update(&changedEntry)
			if reflect.DeepEqual(changedEntry, updatedEntry) || err != nil {
				t.Error("Failed to update entry")
			}
		})

		t.Run("Update entry with fail (entry not found)", func(t *testing.T) {
			unknownEntry := model.UserShoppingListEntry{
				ShoppingListId: 10,
				ProductId:      100,
				Count:          30,
				Note:           "I dont want this anymore",
				Checked:        false,
			}

			_, err := rqliteRepository.Update(&unknownEntry)
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
		entryToDelete := model.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      1,
			Count:          1,
			Note:           "This is very important to me",
			Checked:        false,
		}

		rqliteRepository.Create(&entryToDelete)

		t.Run("Delete entry with success", func(t *testing.T) {
			err = rqliteRepository.Delete(&entryToDelete)
			if err != nil {
				t.Errorf("Failed to delete entry with shoppingListId %d and productId %d", entryToDelete.ShoppingListId, entryToDelete.ProductId)
			}
		})

		t.Run("Delete entry with fail (entry not found)", func(t *testing.T) {
			unknownEntry := model.UserShoppingListEntry{
				ShoppingListId: 10,
				ProductId:      15,
				Count:          155,
				Note:           "This is very important to me",
				Checked:        false,
			}
			err := rqliteRepository.Delete(&unknownEntry)
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
