package handler

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type registerResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  int    `json:"role"`
}

func (r *registerRequest) isValid() bool {
	return r.Email != "" && r.Password != "" && r.Name != ""
}

type RegisterHandler struct {
	userRepository user.Repository
	hasher         crypto.Hasher
}

func NewRegisterHandler(
	userRepository user.Repository,
	hasher crypto.Hasher,
) *RegisterHandler {
	return &RegisterHandler{userRepository, hasher}
}

func (handler *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request registerRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.isValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if model.Role(request.Role) == model.Administrator {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		foundUser, err := handler.userRepository.FindByEmail(request.Email)
		if err != nil && err.Error() != user.ErrorUserNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if foundUser != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		hashedPassword, err := handler.hasher.Hash([]byte(request.Password))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		createdUser, err := handler.userRepository.Create(&model.User{
			Email:    request.Email,
			Password: hashedPassword,
			Name:     request.Name,
			Role:     model.Role(request.Role),
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := registerResponse{
			Id:    createdUser.Id,
			Name:  createdUser.Name,
			Email: createdUser.Email,
			Role:  int(createdUser.Role),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
