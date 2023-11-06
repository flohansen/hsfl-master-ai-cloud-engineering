package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	listModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry"
	entryModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"log"
	"net/http"
)

func main() {
	shoppingListRepository := userShoppingList.NewDemoRepository()
	shoppingListController := userShoppingList.NewDefaultController(shoppingListRepository)
	createContentForShoppingLists(shoppingListRepository)

	shoppingListEntryRepository := userShoppingListEntry.NewDemoRepository()
	shoppingListEntryController := userShoppingListEntry.NewDefaultController(shoppingListEntryRepository)
	createContentForShoppingListEntries(shoppingListEntryRepository)

	handler := router.New(shoppingListController, shoppingListEntryController)

	if err := http.ListenAndServe(":3002", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

func createContentForShoppingListEntries(shoppingListEntryRepository userShoppingListEntry.Repository) {
	shoppingListEntries := []*entryModel.UserShoppingListEntry{
		{
			ShoppingListId: 1,
			ProductId:      1,
			Count:          1,
			Note:           "",
			Checked:        false,
		},
		{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          1,
			Note:           "",
			Checked:        false,
		},
		{
			ShoppingListId: 2,
			ProductId:      3,
			Count:          1,
			Note:           "",
			Checked:        false,
		},
	}

	for _, price := range shoppingListEntries {
		_, err := shoppingListEntryRepository.Create(price)
		if err != nil {
			return
		}
	}
}

func createContentForShoppingLists(shoppingListRepository userShoppingList.Repository) {
	shoppingLists := []*listModel.UserShoppingList{
		{
			Id:          1,
			UserId:      2,
			Description: "Frühstück mit Anne",
			Completed:   false,
		},
		{
			Id:          2,
			UserId:      2,
			Description: "Geburtstagskuchen",
			Completed:   false,
		},
		{
			Id:          3,
			UserId:      4,
			Description: "Einkauf für die Woche",
			Completed:   true,
		},
	}

	for _, price := range shoppingLists {
		_, err := shoppingListRepository.Create(price)
		if err != nil {
			return
		}
	}
}
