package books

import (
	"encoding/json"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"net/http"
	"strconv"
)

type createBookRequest struct {
	Name        string `json:"name"`
	Author      string `json:"author"` //Isnt that for user input, should stuff like author and bookID be included on create?
	Description string `json:"description"`
}

type updateBookRequest struct {
	Name        string `json:"name"`
	Author      string `json:"author"` //Does Update Author make sense?
	Description string `json:"description"`
}

type createChapterRequest struct {
	BookID  int64  `json:"bookid"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	Price   int64  `json:"price"`
	Content string `json:"content"`
}

type updateChapterRequest struct {
	BookID  int64  `json:"bookid"` //Does Update BookID and Author make sense? No, right?
	Name    string `json:"name"`
	Author  string `json:"author"`
	Price   int64  `json:"price"`
	Content string `json:"content"`
}

func (r createBookRequest) isValid() bool {
	return r.Name != "" && r.Author != ""
}

func (r createChapterRequest) isValid() bool {
	return r.Name != "" && r.Author != ""
}

type DefaultController struct {
	bookRepository    BookRepository
	chapterRepository ChapterRepository
}

func NewDefaultController(
	bookRepository BookRepository,
	chapterRepository ChapterRepository,
) *DefaultController {
	return &DefaultController{bookRepository, chapterRepository}
}

func (ctrl *DefaultController) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := ctrl.bookRepository.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (ctrl *DefaultController) PostBooks(w http.ResponseWriter, r *http.Request) {
	var request createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.bookRepository.Create([]*model.Book{{
		Name:        request.Name,
		Author:      request.Author,
		Description: request.Description,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetBook(w http.ResponseWriter, r *http.Request) {
	bookId := r.Context().Value("bookid").(string)

	id, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := ctrl.bookRepository.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (ctrl *DefaultController) PutBook(w http.ResponseWriter, r *http.Request) {
	bookId := r.Context().Value("bookid").(string)

	id, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request updateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.bookRepository.Create([]*model.Book{{
		ID:          id,
		Name:        request.Name,
		Author:      request.Author,
		Description: request.Description,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId := r.Context().Value("bookid").(string)

	id, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.bookRepository.Delete([]*model.Book{{ID: id}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//Chapter Functions

func (ctrl *DefaultController) GetChapters(w http.ResponseWriter, r *http.Request) {
	chapters, err := ctrl.chapterRepository.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapters)
}

func (ctrl *DefaultController) PostChapters(w http.ResponseWriter, r *http.Request) {
	var request createChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.chapterRepository.Create([]*model.Chapter{{
		BookID:  request.BookID,
		Name:    request.Name,
		Author:  request.Author,
		Price:   request.Price,
		Content: request.Content,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetChapter(w http.ResponseWriter, r *http.Request) {
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseInt(chapterId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chapter, err := ctrl.chapterRepository.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapter)
}

func (ctrl *DefaultController) PutChapter(w http.ResponseWriter, r *http.Request) {
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseInt(chapterId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request updateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.chapterRepository.Create([]*model.Chapter{{
		ID:      id,
		BookID:  request.BookID,
		Name:    request.Name,
		Author:  request.Author,
		Price:   request.Price,
		Content: request.Content,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseInt(chapterId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.chapterRepository.Delete([]*model.Chapter{{ID: id}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
