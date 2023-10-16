package transactions

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"

type TransactionRepository interface {
	Migrate() error
	Create([]*model.Transaction) error
	FindAll() ([]*model.Transaction, error)
	FindById(id uint64) (*model.Transaction, error)
	// No delete for now Delete([]*model.Transaction) error
}
