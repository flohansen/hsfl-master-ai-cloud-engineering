package transactions

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
)

type Controller interface {
	GetTransactions(http.ResponseWriter, *http.Request)
	PostTransactions(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	AuthenticationMiddleware(http.ResponseWriter, *http.Request, router.Next)
	//No refund for now, so no delete, nor put
}
