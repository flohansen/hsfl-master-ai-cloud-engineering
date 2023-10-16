package transactions

import (
	"encoding/json"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"net/http"
	"strconv"
)

type createTransactionRequest struct {
	ChapterID       uint64 `json:"chapterID"`
	PayingUserID    uint64 `json:"payingUserID"`
	ReceivingUserID uint64 `json:"receivingUserID"`
	Amount          uint64 `json:"amount"`
}

func (r createTransactionRequest) isValid() bool {
	return r.ChapterID != 0 && r.PayingUserID != 0 && r.ReceivingUserID != 0
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
		ChapterID:       request.ChapterID,
		PayingUserID:    request.PayingUserID,
		ReceivingUserID: request.ReceivingUserID,
		Amount:          request.Amount,
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
