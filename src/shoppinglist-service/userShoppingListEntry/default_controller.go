package userShoppingListEntry

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"net/http"
	"strconv"
)

type defaultController struct {
	userShoppingListEntryRepository Repository
}

func NewDefaultController(userShoppingListEntryRepository Repository) *defaultController {
	return &defaultController{userShoppingListEntryRepository}
}

func (controller defaultController) GetEntries(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	values, err := controller.userShoppingListEntryRepository.FindAll(listId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(values)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (controller defaultController) GetEntry(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := controller.userShoppingListEntryRepository.FindByIds(listId, productId)
	if err != nil {
		if err.Error() == ErrorEntryNotFound {
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
func (controller defaultController) PostEntry(writer http.ResponseWriter, request *http.Request) {
	listId, listIdErr := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	productId, productIdErr := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if listIdErr != nil || productIdErr != nil {
		http.Error(writer, "Invalid listId or productId", http.StatusBadRequest)
		return
	}

	var requestData JsonFormatCreateEntryRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.userShoppingListEntryRepository.Create(&model.UserShoppingListEntry{
		ShoppingListId: listId,
		ProductId:      productId,
		Count:          requestData.Count,
		Note:           requestData.Note,
		Checked:        requestData.Checked,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (controller defaultController) PutEntry(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData JsonFormatCreateEntryRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.userShoppingListEntryRepository.Update(&model.UserShoppingListEntry{
		ShoppingListId: listId,
		ProductId:      productId,
		Count:          requestData.Count,
		Note:           requestData.Note,
		Checked:        requestData.Checked,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (controller defaultController) DeleteEntry(writer http.ResponseWriter, request *http.Request) {
	listId, err := strconv.ParseUint(request.Context().Value("listId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := controller.userShoppingListEntryRepository.Delete(&model.UserShoppingListEntry{ShoppingListId: listId, ProductId: productId}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
