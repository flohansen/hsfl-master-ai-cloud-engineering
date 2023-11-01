package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry"
	"log"
	"net/http"
)

func main() {
	shoppingListRepository := userShoppingList.NewDemoRepository()
	shoppingListController := userShoppingList.NewDefaultController(shoppingListRepository)

	shoppingListEntryRepository := userShoppingListEntry.NewDemoRepository()
	shoppingListEntryController := userShoppingListEntry.NewDefaultController(shoppingListEntryRepository)

	handler := router.New(shoppingListController, shoppingListEntryController)

	if err := http.ListenAndServe(":3002", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
