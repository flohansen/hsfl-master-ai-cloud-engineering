package user

import (
	"encoding/json"
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
	"strconv"
)

type defaultController struct {
	userRepository Repository
}

func NewDefaultController(userRepository Repository) *defaultController {
	return &defaultController{userRepository}
}

func (controller defaultController) GetUser(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := controller.userRepository.FindById(userId)
	if err != nil {
		if err.Error() == ErrorUserNotFound {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (controller defaultController) PostUser(writer http.ResponseWriter, request *http.Request) {
	var requestData JsonFormatCreateUserRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bcryptHasher := crypto.NewBcryptHasher()
	hashedPassword, _ := bcryptHasher.Hash(requestData.Password)
	fmt.Println(requestData)

	if _, err := controller.userRepository.Create(&model.User{
		Email:    requestData.Email,
		Password: hashedPassword,
		Name:     requestData.Name,
		Role:     model.Customer,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}
