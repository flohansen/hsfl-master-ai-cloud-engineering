//go:build integration
// +build integration

package user

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"reflect"
	"testing"
	"time"
)

const TestPort = "7001"

func TestIntegrationRQLiteRepository(t *testing.T) {
	container, err := prepareIntegrationTestRQLiteDatabase()
	if err != nil {
		t.Error(err)
	}

	rqliteRepository := NewRQLiteRepository("http://localhost:" + TestPort + "/?disableClusterDiscovery=true")

	t.Run("TestIntegrationRQLiteRepository_Create", func(t *testing.T) {
		user := model.User{
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		}

		t.Run("Create user with success", func(t *testing.T) {
			_, err = rqliteRepository.Create(&user)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("Can't create user with existing mail address", func(t *testing.T) {
			_, err = rqliteRepository.Create(&user)
			if err.Error() != ErrorUserAlreadyExists {
				t.Error(err)
			}
		})

		t.Run("Database error should return error", func(t *testing.T) {
			_, err = rqliteRepository.Create(&user)
			if err == nil {
				t.Error("there should be an error")
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindAll", func(t *testing.T) {
		users := []*model.User{
			{
				Id:       1,
				Email:    "ada.lovelace@gmail.com",
				Password: []byte("123456"),
				Name:     "Ada Lovelace",
				Role:     model.Customer,
			},
			{
				Id:       2,
				Email:    "alan.turin@gmail.com",
				Password: []byte("123456"),
				Name:     "Alan Turing",
				Role:     model.Customer,
			},
		}

		for _, user := range users {
			_, err := rqliteRepository.Create(user)
			if err != nil {
				t.Error(err)
			}
		}

		t.Run("Successfully fetch all users", func(t *testing.T) {
			fetchedUsers, err := rqliteRepository.FindAll()

			if err != nil {
				t.Error("Can't fetch users")
			}

			if len(fetchedUsers) != len(users) {
				t.Errorf("Unexpected user count. Expected %d, got %d", len(users), len(fetchedUsers))
			}

			if !reflect.DeepEqual(users, fetchedUsers) {
				t.Error("Fetched users do not match expected users")
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindAllByRole", func(t *testing.T) {
		customer := []*model.User{
			{
				Email:    "agathe.bauer@gmail.com",
				Password: []byte("123456"),
				Name:     "Agathe Bauer",
				Role:     model.Customer,
			},
			{
				Email:    "bob.marley@gmail.com",
				Password: []byte("123456"),
				Name:     "Bob Marley",
				Role:     model.Customer,
			},
		}

		merchants := []*model.User{
			{
				Email:    "ada.lovelace@gmail.com",
				Password: []byte("123456"),
				Name:     "Ada Lovelace",
				Role:     model.Merchant,
			},
			{
				Email:    "alan.turin@gmail.com",
				Password: []byte("123456"),
				Name:     "Alan Turing",
				Role:     model.Merchant,
			},
		}

		for _, user := range customer {
			_, err := rqliteRepository.Create(user)
			if err != nil {
				t.Error(err)
			}
		}

		for _, user := range merchants {
			_, err := rqliteRepository.Create(user)
			if err != nil {
				t.Error(err)
			}
		}

		t.Run("Succesfully fetch all merchants", func(t *testing.T) {
			fetchedUsers, err := rqliteRepository.FindAllByRole(model.Merchant)

			if err != nil {
				t.Error("Can't fetch users")
			}

			if len(fetchedUsers) != len(merchants) {
				t.Errorf("Unexpected user count. Expected %d, got %d", len(merchants), len(fetchedUsers))
			}

			if !reflect.DeepEqual(merchants, fetchedUsers) {
				t.Error("Fetched users do not match expected users")
			}
		})

		t.Run("Succesfully fetch all customers", func(t *testing.T) {
			fetchedUsers, err := rqliteRepository.FindAllByRole(model.Customer)

			if err != nil {
				t.Error("Can't fetch users")
			}

			if len(fetchedUsers) != len(merchants) {
				t.Errorf("Unexpected user count. Expected %d, got %d", len(customer), len(fetchedUsers))
			}

			if !reflect.DeepEqual(customer, fetchedUsers) {
				t.Error("Fetched users do not match expected users")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindById", func(t *testing.T) {
		user := model.User{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		}

		rqliteRepository.Create(&user)

		t.Run("Successfully fetch user", func(t *testing.T) {
			fetchedUser, err := rqliteRepository.FindById(user.Id)
			if err != nil {
				t.Errorf("Can't find expected user with id %d: %v", user.Id, err)
			}

			if !reflect.DeepEqual(user, *fetchedUser) {
				t.Error("Fetched user does not match original user")
			}
		})

		t.Run("Fail to fetch user (user not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindById(2)
			if err == nil {
				t.Errorf("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindByEmail", func(t *testing.T) {
		user := model.User{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		}

		rqliteRepository.Create(&user)

		t.Run("Successfully fetch user", func(t *testing.T) {
			fetchedUser, err := rqliteRepository.FindByEmail(user.Email)
			if err != nil {
				t.Errorf("Can't find expected user with id %d: %v", user.Id, err)
			}

			if !reflect.DeepEqual(user, *fetchedUser) {
				t.Error("Fetched user does not match original user")
			}
		})

		t.Run("Fail to fetch user (user not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindByEmail("unknown@mail.com")
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
		user := model.User{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		}

		rqliteRepository.Create(&user)

		t.Run("Update user with success", func(t *testing.T) {
			changedUser := model.User{
				Id:       1,
				Email:    "ada.lovelace@gmail.com",
				Password: []byte("123456"),
				Name:     "Ada Lovelace",
				Role:     model.Customer,
			}

			updatedUser, err := rqliteRepository.Update(&changedUser)
			if updatedUser.Name != updatedUser.Name || err != nil {
				t.Error("Failed to update user")
			}
		})

		t.Run("Update user with fail (user not found)", func(t *testing.T) {
			changedUser := model.User{
				Id:       2,
				Email:    "bob.marley@gmail.com",
				Password: []byte("123456"),
				Name:     "Bob Marley",
				Role:     model.Customer,
			}

			_, err := rqliteRepository.Update(&changedUser)
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
		userToDelete := model.User{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		}

		rqliteRepository.Create(&userToDelete)

		t.Run("Delete user with success", func(t *testing.T) {
			err = rqliteRepository.Delete(&userToDelete)
			if err != nil {
				t.Errorf("Failed to delete user with id %d", userToDelete.Id)
			}
		})

		t.Run("Delete user with fail (user not found)", func(t *testing.T) {
			nonExistingUser := model.User{
				Id:       2,
				Email:    "bob.marley@gmail.com",
				Password: []byte("123456"),
				Name:     "Bob Marley",
				Role:     model.Customer,
			}

			err := rqliteRepository.Delete(&nonExistingUser)
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
