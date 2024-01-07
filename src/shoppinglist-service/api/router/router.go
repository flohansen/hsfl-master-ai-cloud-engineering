package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	shoppingListController *userShoppingList.Controller,
	shoppingListEntryController *userShoppingListEntry.Controller,
	authMiddleware router.Middleware,
) *Router {
	r := router.New()

	r.GET("/api/v1/shoppinglist/:userId", (*shoppingListController).GetLists, authMiddleware)
	r.GET("/api/v1/shoppinglist/:listId/:userId", (*shoppingListController).GetList, authMiddleware)
	r.PUT("/api/v1/shoppinglist/:listId/:userId", (*shoppingListController).PutList, authMiddleware)
	r.POST("/api/v1/shoppinglist/:userId", (*shoppingListController).PostList, authMiddleware)
	r.DELETE("/api/v1/shoppinglist/:listId", (*shoppingListController).DeleteList, authMiddleware)

	r.GET("/api/v1/shoppinglistentries/:listId", (*shoppingListEntryController).GetEntries, authMiddleware)
	r.GET("/api/v1/shoppinglistentries/:listId/:productId", (*shoppingListEntryController).GetEntry, authMiddleware)
	r.PUT("/api/v1/shoppinglistentries/:listId/:productId", (*shoppingListEntryController).PutEntry, authMiddleware)
	r.POST("/api/v1/shoppinglistentries/:listId/:productId", (*shoppingListEntryController).PostEntry, authMiddleware)
	r.DELETE("/api/v1/shoppinglistentries/:listId/:productId", (*shoppingListEntryController).DeleteEntry, authMiddleware)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
