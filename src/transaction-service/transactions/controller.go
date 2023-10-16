package transactions

import "net/http"

type Controller interface {
	GetTransactions(http.ResponseWriter, *http.Request)
	PostTransactions(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	//No refund for now, so no delete, nor put
}
