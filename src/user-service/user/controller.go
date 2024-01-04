package user

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
)

type JsonFormatUpdateUserRequest struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

type JsonFormatCreateUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password []byte `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}

type JsonFormatGetUserResponse struct {
	ID    uint64     `json:"id,omitempty"`
	Email string     `json:"email,omitempty"`
	Name  string     `json:"name,omitempty"`
	Role  model.Role `json:"role,omitempty"`
}

type JsonFormatGetUserByRoleResponse struct {
	ID   uint64     `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`
	Role model.Role `json:"role,omitempty"`
}

type Controller interface {
	GetUsersByRole(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	PutUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
}
