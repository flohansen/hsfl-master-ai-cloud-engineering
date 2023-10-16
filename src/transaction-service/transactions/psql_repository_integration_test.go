package transactions

import (
	"context"
	"database/sql"
	bookModels "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPsqlRepository(t *testing.T) {
	postgres, err := containerhelpers.StartPostgres()
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err.Error())
	}

	t.Cleanup(func() {
		postgres.Terminate(context.Background())
	})

	port, err := postgres.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Fatalf("could not get database container port: %s", err.Error())
	}

	repository, err := NewPsqlRepository(database.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	if err != nil {
		t.Fatalf("could not create transaction repository: %s", err.Error())
	}
	t.Cleanup(clearTables(t, repository.db))

	t.Run("Migrate", func(t *testing.T) {
		t.Run("should create transactions table", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			createUserBookAndChapterTable(t, repository.db)
			// when
			err := repository.Migrate()

			// then
			assert.NoError(t, err)
			assertTableExists(t, repository.db, "transactions", []string{"id"})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("should insert users in batches", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given

			transactions := []*model.Transaction{
				{
					ID:              1,
					ChapterID:       1,
					PayingUserID:    1,
					ReceivingUserID: 1,
					Amount:          0,
				},
				{
					ID:              2,
					ChapterID:       1,
					PayingUserID:    1,
					ReceivingUserID: 1,
					Amount:          0,
				},
			}
			// when
			err := repository.Create(transactions)

			// then
			assert.NoError(t, err)
			assert.Equal(t, transactions[0], getTransactionFromDatabase(t, repository.db, 1))
			assert.Equal(t, transactions[1], getTransactionFromDatabase(t, repository.db, 2))
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("should return transaction", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			insertTransaction(t, repository.db, &model.Transaction{
				ID:              3,
				ChapterID:       1,
				PayingUserID:    1,
				ReceivingUserID: 1,
				Amount:          0,
			})

			// when
			transaction, err := repository.FindByID(3)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, transaction)
		})
	})
}

func createUserBookAndChapterTable(t *testing.T, db *sql.DB) {
	db.Exec(`
		create table if not exists users (
			id				serial primary key, 
			email			varchar(100) not null unique,
			username    	varchar(16) not null unique,
			password 		bytea not null,
			profile_name 	varchar(100) not null,
			balance 		int not null default 0					
		)
	`)

	db.Exec(`
		create table if not exists books (
			id			serial primary key,
			name    	varchar(100) not null,
			authorId	int not null,
			description text not null,
			foreign key (authorId) REFERENCES users(id)
		)
	`)

	db.Exec(`
		create table if not exists chapters (
    		id			serial primary key,
    		bookId		int not null,
			name    	varchar(100) not null,
			authorId	int not null,
			price		int not null,
			content 	text not null,
   			foreign key (bookId) REFERENCES books(id),
   			foreign key (authorId) REFERENCES users(id)
		)
	`)

	db.Exec(`insert into users (email, username, password, profile_name) values ($1,$2,$3,$4)`, "1a@mail.com", "tester1", []byte("pw"), "Peter")
	db.Exec(`insert into users (email, username, password, profile_name) values ($1,$2,$3,$4)`, "2a@mail.com", "tester2", []byte("pw"), "Ursula")

	books := []*bookModels.Book{
		{
			ID:          1,
			Name:        "Book One",
			AuthorID:    1,
			Description: "A good book",
		},
		{
			ID:          2,
			Name:        "Book Two",
			AuthorID:    1,
			Description: "A bad book",
		},
	}
	for _, book := range books {
		insertBook(t, db, book)
	}

	chapters := []*bookModels.Chapter{
		{
			ID:       1,
			BookID:   1,
			Name:     "doesnt matter",
			AuthorID: 1,
			Price:    0,
			Content:  "doesnt matter",
		},
		{
			ID:       2,
			BookID:   1,
			Name:     "doesnt matter",
			AuthorID: 1,
			Price:    0,
			Content:  "doesnt matter",
		},
	}
	for _, chapter := range chapters {
		insertChapter(t, db, chapter)
	}

}

func getTransactionFromDatabase(t *testing.T, db *sql.DB, id int) *model.Transaction {
	row := db.QueryRow(`select id, chapterid, payinguserid, receivinguserid, amount from transactions where id = $1`, id)

	var transaction model.Transaction
	if err := row.Scan(&transaction.ID, &transaction.ChapterID, &transaction.PayingUserID, &transaction.ReceivingUserID, &transaction.Amount); err != nil {
		return nil
	}

	return &transaction
}

func insertTransaction(t *testing.T, db *sql.DB, transaction *model.Transaction) {
	_, err := db.Exec(`insert into transactions (id, chapterid, payinguserid, receivinguserid, amount) values ($1, $2, $3, $4, $5)`, transaction.ID, transaction.ChapterID, transaction.PayingUserID, transaction.ReceivingUserID, transaction.Amount)
	if err != nil {
		t.Logf("could not insert transaction: %s", err.Error())
		t.FailNow()
	}
}
func insertBook(t *testing.T, db *sql.DB, book *bookModels.Book) {
	_, err := db.Exec(`insert into books (name, authorId, description) values ($1, $2, $3)`, book.Name, book.AuthorID, book.Description)
	if err != nil {
		t.Logf("could not insert book: %s", err.Error())
		t.FailNow()
	}
}
func insertChapter(t *testing.T, db *sql.DB, chapter *bookModels.Chapter) {
	_, err := db.Exec(`insert into chapters (id,bookId,name,authorId,price,content) values ($1, $2, $3, $4, $5, $6)`, chapter.ID, chapter.BookID, chapter.Name, chapter.AuthorID, chapter.Price, chapter.Content)
	if err != nil {
		t.Logf("could not insert chapter: %s", err.Error())
		t.FailNow()
	}
}

func clearTables(t *testing.T, db *sql.DB) func() {
	return func() {
		if _, err := db.Exec("delete from transactions"); err != nil {
			t.Logf("could not delete rows from transactions: %s", err.Error())
			t.FailNow()
		}
	}
}

func assertTableExists(t *testing.T, db *sql.DB, name string, columns []string) {
	rows, err := db.Query(`select column_name from information_schema.columns where table_name = $1`, name)
	if err != nil {
		t.Fail()
		return
	}

	scannedCols := make(map[string]struct{})
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			t.Logf("expected")
			t.FailNow()
		}

		scannedCols[column] = struct{}{}
	}

	if len(scannedCols) == 0 {
		t.Logf("expected table '%s' to exist, but not found", name)
		t.FailNow()
	}

	for _, col := range columns {
		if _, ok := scannedCols[col]; !ok {
			t.Logf("expected table '%s' to have column '%s'", name, col)
			t.Fail()
		}
	}
}
