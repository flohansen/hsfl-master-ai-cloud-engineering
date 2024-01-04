package user

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"reflect"
	"sort"
	"testing"
)

func TestDemoRepository_CreateUser(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	user := model.User{
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	// Create user with success
	_, err := demoRepository.Create(&user)
	if err != nil {
		t.Error(err)
	}

	// Check for doublet
	_, err = demoRepository.Create(&user)
	if err.Error() != ErrorUserAlreadyExists {
		t.Error(err)
	}
}

func TestDemoRepository_FindAll(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

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
		_, err := demoRepository.Create(user)
		if err != nil {
			t.Error("Failed to add prepared product for test")
		}
	}

	t.Run("Fetch all users", func(t *testing.T) {
		fetchedUsers, err := demoRepository.FindAll()
		if err != nil {
			t.Error("Can't fetch users")
		}

		if len(fetchedUsers) != len(users) {
			t.Errorf("Unexpected product count. Expected %d, got %d", len(users), len(fetchedUsers))
		}
	})

	userTests := []struct {
		name string
		want *model.User
	}{
		{
			name: "first",
			want: users[0],
		},
		{
			name: "second",
			want: users[1],
		},
	}

	for i, tt := range userTests {
		t.Run("Is fetched product matching with "+tt.name+" added product?", func(t *testing.T) {
			fetchedUsers, _ := demoRepository.FindAll()
			sort.Slice(fetchedUsers, func(i, j int) bool {
				return fetchedUsers[i].Id < fetchedUsers[j].Id
			})
			if !reflect.DeepEqual(tt.want, fetchedUsers[i]) {
				t.Error("Fetched product does not match original product")
			}
		})
	}
}

func TestDemoRepository_FindAllByRole(t *testing.T) {
	demoRepository := NewDemoRepository()

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
			Role:     model.Merchant,
		},
	}

	for _, user := range users {
		_, err := demoRepository.Create(user)
		if err != nil {
			t.Error("Failed to add prepared user for test")
		}
	}

	t.Run("Fetch all users by merchant role", func(t *testing.T) {
		merchantRole := model.Merchant
		fetchedUsers, err := demoRepository.FindAllByRole(&merchantRole)
		if err != nil {
			t.Error("Can't fetch users")
		}

		if len(fetchedUsers) != 1 {
			t.Errorf("Unexpected user count. Expected 1, got %d", len(fetchedUsers))
		}
	})
}

func TestDemoRepository_FindById(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	user := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	_, err := demoRepository.Create(&user)
	if err != nil {
		t.Error("Failed to add prepare user for test")
	}

	// Fetch user with existing id
	fetchedUser, err := demoRepository.FindById(user.Id)
	if err != nil {
		t.Errorf("Can't find expected user with id %d", user.Id)
	}

	// Is fetched user matching with added user?
	if !reflect.DeepEqual(user, *fetchedUser) {
		t.Error("Fetched user does not match original user")
	}

	// Non-existing user test
	_, err = demoRepository.FindById(42)
	if err.Error() != "user could not be found" {
		t.Error(err)
	}
}

func TestDemoRepository_FindByEmail(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	user := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	_, err := demoRepository.Create(&user)
	if err != nil {
		t.Error("Failed to add prepare user for test")
	}

	// Fetch user with existing email
	fetchedUser, err := demoRepository.FindByEmail(user.Email)
	if err != nil {
		t.Errorf("Can't find expected user with email %s", user.Email)
	}

	// Is fetched user matching with added user?
	if !reflect.DeepEqual(user, *fetchedUser) {
		t.Error("Fetched user does not match original user")
	}

	// Non-existing user test
	_, err = demoRepository.FindByEmail("alan.turing@gmail.com")
	if err.Error() != "user could not be found" {
		t.Error(err)
	}
}

func TestDemoRepository_Update(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	user := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	fetchedUser, err := demoRepository.Create(&user)
	if err != nil {
		t.Error("Failed to add prepare user for test")
	}

	fetchedUser.Name = "Alan Turing"
	updatedUser, err := demoRepository.Update(fetchedUser)

	// Check if returned user has the updated name
	if fetchedUser.Name != updatedUser.Name {
		t.Error("Failed to update user")
	}
}

func TestDemoRepository_Delete(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	user := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	fetchedUser, err := demoRepository.Create(&user)
	if err != nil {
		t.Error("Failed to add prepare user for test")
	}

	// Test for deletion
	err = demoRepository.Delete(fetchedUser)
	if err != nil {
		t.Errorf("Failed to delete user with id %d", user.Id)
	}

	// Fetch user with existing id
	fetchedUser, err = demoRepository.FindById(user.Id)
	if err.Error() != "user could not be found" {
		t.Errorf("User with id %d was not deleted", user.Id)
	}

	// Try to delete non-existing user
	fakeUser := model.User{
		Id:       42,
		Email:    "alan.turin@gmail.com",
		Password: []byte("123456"),
		Name:     "Alan Turing",
		Role:     model.Customer,
	}

	// Test for deletion
	err = demoRepository.Delete(&fakeUser)
	if err.Error() != "user could not be deleted" {
		t.Errorf("User with id %d was deleted", user.Id)
	}
}
