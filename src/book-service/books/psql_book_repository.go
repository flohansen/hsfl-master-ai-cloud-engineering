package books

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	_ "github.com/lib/pq"
)

type PsqlBookRepository struct {
	db *sql.DB
}

func NewPsqlBookRepository(config database.Config) (*PsqlBookRepository, error) {
	dsn := config.Dsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PsqlBookRepository{db}, nil
}

const createBooksTable = `
create table if not exists books (
    id			serial primary key,
	name    	varchar(100) not null,
	authorId	int not null,
	description text not null,
   	foreign key (authorId) REFERENCES users(id)
)
`

func (repo *PsqlBookRepository) Migrate() error {
	_, err := repo.db.Exec(createBooksTable)
	return err
}

const createBooksBatchQuery = `
insert into books (name, authorId, description) values %s
`

func (repo *PsqlBookRepository) Create(books []*model.Book) error {
	placeholders := make([]string, len(books))
	values := make([]interface{}, len(books)*3)

	for i := 0; i < len(books); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d)", i*3+1, i*3+2, i*3+3)
		values[i*3+0] = books[i].Name
		values[i*3+1] = books[i].AuthorID
		values[i*3+2] = books[i].Description
	}

	query := fmt.Sprintf(createBooksBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const updateBookBatchQuery = `
update books set name = $1, description = $2 where id = $3
`

func (repo *PsqlBookRepository) Update(id uint64, book *model.UpdateBook) error {
	_, err := repo.db.Exec(updateBookBatchQuery, book.Name, book.Description, id)
	return err
}

const findAllBooksQuery = `
select id, name, authorId, description from books
`

func (repo *PsqlBookRepository) FindAll() ([]*model.Book, error) {
	rows, err := repo.db.Query(findAllBooksQuery)
	if err != nil {
		return nil, err
	}

	books := make([]*model.Book, 0)
	for rows.Next() {
		book := model.Book{}
		if err := rows.Scan(&book.ID, &book.Name, &book.AuthorID, &book.Description); err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	return books, nil
}

const findBooksByIDQuery = `
select id, name, authorId, description from books where id = $1 limit 1
`

func (repo *PsqlBookRepository) FindById(id uint64) (*model.Book, error) {
	row := repo.db.QueryRow(findBooksByIDQuery, id)

	var book model.Book
	if err := row.Scan(&book.ID, &book.Name, &book.AuthorID, &book.Description); err != nil {
		return nil, err
	}
	return &book, nil
}

const deleteBooksBatchQuery = `
delete from books where id in (%s)
`

func (repo *PsqlBookRepository) Delete(books []*model.Book) error {
	placeholders := make([]string, len(books))
	ids := make([]interface{}, len(books))

	for i := 0; i < len(books); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		ids[i] = books[i].ID
	}

	query := fmt.Sprintf(deleteBooksBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, ids...)
	return err
}
