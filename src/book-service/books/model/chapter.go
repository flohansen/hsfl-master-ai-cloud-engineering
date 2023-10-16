package model

type UpdateChapter struct {
	Name    string `json:"name"`
	Price   uint64 `json:"price"`
	Content string `json:"content"`
}
type Chapter struct {
	ID       uint64 `json:"id"`
	BookID   uint64 `json:"bookid"`
	Name     string `json:"name"`
	AuthorID uint64 `json:"authorid"`
	Price    uint64 `json:"price"`
	Content  string `json:"content"`
}
