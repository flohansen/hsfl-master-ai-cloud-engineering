package user

import (
	"encoding/json"
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

	_, exists := request.Context().Value("auth_userId").(uint64)
	if !exists {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	values, err := controller.userRepository.FindAllByRole(model.Role(userRole))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(values) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	formattedResponse := make([]*JsonFormatGetUserByRoleResponse, len(values))
	for i, value := range values {
		formattedResponse[i] = &JsonFormatGetUserByRoleResponse{
			ID:   value.Id,
			Name: value.Name,
			Role: value.Role,
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(formattedResponse)
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

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(model.Role)

	value, err := controller.userRepository.FindById(userId)

	if (authUserId == userId) ||
		(authUserRole == model.Administrator) ||
		(value != nil && value.Role == model.Merchant) {

		if err != nil {
			if err.Error() == ErrorUserNotFound {
				http.Error(writer, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		responseValue := JsonFormatGetUserResponse{
			ID:    value.Id,
			Email: value.Email,
			Name:  value.Name,
			Role:  value.Role,
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(responseValue)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
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

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(model.Role)

	if (authUserId == userId) ||
		(authUserRole == model.Administrator) {
		if _, err := controller.userRepository.Update(&model.User{
			Id:    userId,
			Email: requestData.Email,
			Name:  requestData.Name,
		}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller defaultController) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(model.Role)

	if (authUserId == userId) ||
		(authUserRole == model.Administrator) {
		if err := controller.userRepository.Delete(&model.User{Id: userId}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}
