package transactions

import "net/http"

type Controller interface {
	GetTransactions(http.ResponseWriter, *http.Request)
}
