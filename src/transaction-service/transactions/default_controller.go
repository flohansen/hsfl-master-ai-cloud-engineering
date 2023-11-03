package transactions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"io"
	"net/http"
	"strconv"
)

type createTransactionRequest struct {
	ChapterID    uint64 `json:"chapterID"`
	PayingUserID uint64 `json:"payingUserID"`
	Amount       uint64 `json:"amount"`
}

func (r createTransactionRequest) isValid() bool {
	return r.ChapterID != 0 && r.PayingUserID != 0
}

type DefaultController struct {
	transactionRepository TransactionRepository
}

func NewDefaultController(
	transactionRepository TransactionRepository,
) *DefaultController {
	return &DefaultController{transactionRepository}
}

func (ctrl *DefaultController) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := ctrl.transactionRepository.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (ctrl *DefaultController) PostTransactions(w http.ResponseWriter, r *http.Request) {
	var request createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.transactionRepository.Create([]*model.Transaction{{
		ChapterID:    request.ChapterID,
		PayingUserID: request.PayingUserID,
		Amount:       request.Amount,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetTransaction(w http.ResponseWriter, r *http.Request) {
	transactionId := r.Context().Value("transactionid").(string)

	id, err := strconv.ParseUint(transactionId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transaction, err := ctrl.transactionRepository.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (ctrl *DefaultController) AuthenticationMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/users/me", nil)
	if err != nil {
		http.Error(w, "Response could not be sent!", 500)
		return
	}
	req.Header.Add("Authorization", r.Header.Get("Authorization"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(res.Status, res.Header)

	if res.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var claims map[string]interface{}
	err = json.Unmarshal(bytes, &claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), "user", claims)
	next(r.WithContext(ctx))
}
