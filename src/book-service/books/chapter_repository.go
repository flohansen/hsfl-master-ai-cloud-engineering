package books

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"

type ChapterRepository interface {
	Migrate() error
	Create([]*model.Chapter) error
	Update(id uint64, updateChapter *model.UpdateChapter) error
	FindAll() ([]*model.Chapter, error)
	FindById(id uint64) (*model.Chapter, error)
	Delete([]*model.Chapter) error
}
