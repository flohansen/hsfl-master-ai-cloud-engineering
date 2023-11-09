package user

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"reflect"
	"testing"
)

func TestDemoRepository_CreateUser(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	user := model.User{
		Id:       1,
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
	if err.Error() != "user already exists" {
		t.Error(err)
	}
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
