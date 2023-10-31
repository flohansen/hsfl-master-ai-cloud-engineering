package books

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	_ "github.com/lib/pq"
)

type PsqlChapterRepository struct {
	db *sql.DB
}

func NewPsqlChapterRepository(config database.Config) (*PsqlChapterRepository, error) {
	dsn := config.Dsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PsqlChapterRepository{db}, nil
}

const createChaptersTable = `
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
`

func (repo *PsqlChapterRepository) Migrate() error {
	_, err := repo.db.Exec(createChaptersTable)
	return err
}

const createChaptersBatchQuery = `
insert into chapters (bookId,  name,  authorId,  price,  content) values %s
`

func (repo *PsqlChapterRepository) Create(chapters []*model.Chapter) error {
	placeholders := make([]string, len(chapters))
	values := make([]interface{}, len(chapters)*5)

	for i := 0; i < len(chapters); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		values[i*5+0] = chapters[i].BookID
		values[i*5+1] = chapters[i].Name
		values[i*5+2] = chapters[i].AuthorID
		values[i*5+3] = chapters[i].Price
		values[i*5+4] = chapters[i].Content
	}

	query := fmt.Sprintf(createChaptersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const updateChapterBatchQuery = `
update chapters set name = $1, price = $2, content = $3 where id = $4
`

func (repo *PsqlChapterRepository) Update(id uint64, chapter *model.UpdateChapter) error {
	_, err := repo.db.Exec(updateChapterBatchQuery, chapter.Name, chapter.Price, chapter.Content, id)
	return err
}

const findAllChaptersQuery = `
select id, bookId, name, a  uthorId,price,content from chapters
`

func (repo *PsqlChapterRepository) FindAll() ([]*model.Chapter, error) {
	rows, err := repo.db.Query(findAllChaptersQuery)
	if err != nil {
		return nil, err
	}

	chapters := make([]*model.Chapter, 0)
	for rows.Next() {
		chapter := model.Chapter{}
		if err := rows.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.AuthorID, &chapter.Price, &chapter.Content); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, nil
}

const findChaptersByIDQuery = `
select id, bookId, name, authorId, price, content from chapters where id = $1
`

func (repo *PsqlChapterRepository) FindById(id uint64) (*model.Chapter, error) {
	row := repo.db.QueryRow(findChaptersByIDQuery, id)

	var chapter model.Chapter
	if err := row.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.AuthorID, &chapter.Price, &chapter.Content); err != nil {
		return nil, err
	}

	return &chapter, nil
}

const deleteChaptersBatchQuery = `
delete from chapters where id in (%s)
`

func (repo *PsqlChapterRepository) Delete(chapters []*model.Chapter) error {
	placeholders := make([]string, len(chapters))
	ids := make([]interface{}, len(chapters))

	for i := 0; i < len(chapters); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		ids[i] = chapters[i].ID
	}

	query := fmt.Sprintf(deleteChaptersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, ids...)
	return err
}
