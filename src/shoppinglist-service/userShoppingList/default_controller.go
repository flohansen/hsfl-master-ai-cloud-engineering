package userShoppingList

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"net/http"
	"strconv"
)

type defaultController struct {
	userShoppingListRepository Repository
}

func NewDefaultController(userShoppingListRepository Repository) *defaultController {
	return &defaultController{userShoppingListRepository}
}

func (controller defaultController) GetList(writer http.ResponseWriter, request *http.Request) {
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
}

func (controller defaultController) PutList(writer http.ResponseWriter, request *http.Request) {
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

	if _, err := controller.userShoppingListRepository.Update(&model.UserShoppingList{
		Id:          listId,
		UserId:      userId,
		Description: requestData.Description,
		Completed:   requestData.Checked,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (controller defaultController) PostList(writer http.ResponseWriter, request *http.Request) {
	var requestData JsonFormatCreateListRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.userShoppingListRepository.Create(&model.UserShoppingList{
		UserId:      requestData.UserId,
		Description: requestData.Description,
		Completed:   false,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (controller defaultController) GetLists(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
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
}

func (controller defaultController) DeleteList(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := controller.userShoppingListRepository.Delete(&model.UserShoppingList{Id: listId}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
