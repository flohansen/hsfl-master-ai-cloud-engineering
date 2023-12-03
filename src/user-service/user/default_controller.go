package user

import (
	"encoding/json"
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

func (controller defaultController) GetUsers(writer http.ResponseWriter, request *http.Request) {
	values, err := controller.userRepository.FindAll()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(values)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (controller defaultController) GetUsersByRole(writer http.ResponseWriter, request *http.Request) {
	userRoleStr, ok := request.Context().Value("userRole").(string)
	if !ok {
		http.Error(writer, "Invalid userRole value", http.StatusBadRequest)
		return
	}

	userRole, err := strconv.ParseUint(userRoleStr, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid userRole value", http.StatusBadRequest)
		return
	}

	userRoleModel := getUserRole(userRole)

	if userRoleModel == nil {
		http.Error(writer, "Unknown user role", http.StatusNotFound)
		return
	}

	values, err := controller.userRepository.FindAllByRole(userRoleModel)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(values)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
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

func (controller defaultController) PutUser(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData JsonFormatUpdateUserRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.userRepository.Update(&model.User{
		Id:    userId,
		Email: requestData.Email,
		Name:  requestData.Name,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
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

func (controller defaultController) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := controller.userRepository.Delete(&model.User{Id: userId}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getUserRole(role uint64) *model.Role {
	switch model.Role(role) {
	case model.Customer, model.Merchant, model.Administrator:
		validRole := model.Role(role)
		return &validRole
	default:
		return nil
	}
}
