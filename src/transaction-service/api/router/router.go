package router

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	transactionController transactions.Controller,
) *Router {
	transactionsRouter := router.New()

	transactionsRouter.GET("/api/v1/transactions", transactionController.GetTransactions)
	transactionsRouter.POST("/api/v1/transactions", transactionController.PostTransactions)
	transactionsRouter.GET("/api/v1/transactions/:transactionid", transactionController.GetTransaction)
	return &Router{transactionsRouter}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
