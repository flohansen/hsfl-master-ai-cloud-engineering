package user

import (
	"encoding/json"
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
