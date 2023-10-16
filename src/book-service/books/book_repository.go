package books

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"

type BookRepository interface {
	Migrate() error
	Create([]*model.Book) error
	Update(id uint64, updateBook *model.UpdateBook) error
	FindAll() ([]*model.Book, error)
	FindById(id uint64) (*model.Book, error)
	Delete([]*model.Book) error
}
