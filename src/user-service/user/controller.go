package user

import "net/http"

type Controller interface {
	GetUser(http.ResponseWriter, *http.Request)
}
