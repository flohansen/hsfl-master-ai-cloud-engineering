package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry"
	"log"
	"net/http"
)

func main() {
	shoppingListRepository := userShoppingList.NewDemoRepository()
	shoppingListController := userShoppingList.NewDefaultController(shoppingListRepository)
	createContent(shoppingListRepository)

	shoppingListEntryRepository := userShoppingListEntry.NewDemoRepository()
	shoppingListEntryController := userShoppingListEntry.NewDefaultController(shoppingListEntryRepository)

	handler := router.New(shoppingListController, shoppingListEntryController)

	if err := http.ListenAndServe(":3002", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

func createContent(shoppingListRepository userShoppingList.Repository) {
	shoppingLists := []*model.UserShoppingList{
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
