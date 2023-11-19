package user

import "net/http"

type JsonFormatCreateUserRequest struct {
	Email    string
	Password []byte
	Name     string
}

type Controller interface {
	GetUser(http.ResponseWriter, *http.Request)
	PostUser(http.ResponseWriter, *http.Request)
}
