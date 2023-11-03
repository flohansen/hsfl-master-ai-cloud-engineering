package transactions

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	_ "github.com/lib/pq"
)

type PsqlRepository struct {
	db *sql.DB
}

func NewPsqlRepository(config database.Config) (*PsqlRepository, error) {
	dsn := config.Dsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PsqlRepository{db}, nil
}

const createTransactionsTable = `
create table if not exists transactions (
	id					serial primary key,
	chapterid    		int not null,
	payinguserid 		int not null,
	amount 				int not null,
	foreign key (chapterid) references chapters(id),
	foreign key (payinguserid) references users(id)
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createTransactionsTable)
	return err
}

const createTransactionsBatchQuery = `
insert into transactions (chapterid, payinguserid, amount) values %s
`

func (repo *PsqlRepository) Create(transactions []*model.Transaction) error {
	placeholders := make([]string, len(transactions))
	values := make([]interface{}, len(transactions)*3)

	for i := 0; i < len(transactions); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d)", i*3+1, i*3+2, i*3+3)
		values[i*3+0] = transactions[i].ChapterID
		values[i*3+1] = transactions[i].PayingUserID
		values[i*3+2] = transactions[i].Amount
	}

	query := fmt.Sprintf(createTransactionsBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const findAllTransactionsQuery = `
select id, chapterid, payinguserid, amount from transactions
`

func (repo *PsqlRepository) FindAll() ([]*model.Transaction, error) {
	rows, err := repo.db.Query(findAllTransactionsQuery)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, 0)
	for rows.Next() {
		transaction := model.Transaction{}
		if err := rows.Scan(&transaction.ID, &transaction.ChapterID, &transaction.PayingUserID, &transaction.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

const findTransactionbyIDQuery = `
select id, chapterid, payinguserid, amount from transactions where id = $1
`

func (repo *PsqlRepository) FindById(id uint64) (*model.Transaction, error) {
	row := repo.db.QueryRow(findTransactionbyIDQuery, id)
	transaction := &model.Transaction{}
	if err := row.Scan(&transaction.ID, &transaction.ChapterID, &transaction.PayingUserID, &transaction.Amount); err != nil {
		return nil, err
	}

	return transaction, nil
}
