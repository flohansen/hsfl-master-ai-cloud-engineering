package books

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"

type BookRepository interface {
	Migrate() error
	Create([]*model.Book) error
	FindAll() ([]*model.Book, error)
	FindById(id int64) (*model.Book, error)
	Delete([]*model.Book) error
}
