package model

type UpdateBook struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Book struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	AuthorID    uint64 `json:"authorid"`
	Description string `json:"description"`
}
