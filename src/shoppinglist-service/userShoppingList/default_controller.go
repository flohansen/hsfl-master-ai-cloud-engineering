package userShoppingList

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"net/http"
	"strconv"
)

type DefaultController struct {
	userShoppingListRepository Repository
}

func NewDefaultController(userShoppingListRepository Repository) *DefaultController {
	return &DefaultController{userShoppingListRepository}
}

func (controller DefaultController) GetLists(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserId == userId || authUserRole == auth.Administrator {
		values, err := controller.userShoppingListRepository.FindAllById(userId)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(values)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller DefaultController) GetList(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserId == userId || authUserRole == auth.Administrator {
		value, err := controller.userShoppingListRepository.FindByIds(userId, listId)
		if err != nil {
			if err.Error() == ErrorListNotFound {
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
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller DefaultController) PutList(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData JsonFormatUpdateListRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserId == userId || authUserRole == auth.Administrator {
		if _, err := controller.userShoppingListRepository.Update(&model.UserShoppingList{
			Id:          listId,
			UserId:      userId,
			Description: requestData.Description,
			Completed:   requestData.Checked,
		}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller DefaultController) PostList(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData JsonFormatCreateListRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserId == userId || authUserRole == auth.Administrator {
		if _, err := controller.userShoppingListRepository.Create(&model.UserShoppingList{
			UserId:      userId,
			Description: requestData.Description,
			Completed:   false,
		}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusCreated)
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller DefaultController) DeleteList(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := controller.userShoppingListRepository.FindById(listId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	authUserId, _ := request.Context().Value("auth_userId").(uint64)
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserId == list.UserId || authUserRole == auth.Administrator {
		if err := controller.userShoppingListRepository.Delete(&model.UserShoppingList{Id: listId}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}
