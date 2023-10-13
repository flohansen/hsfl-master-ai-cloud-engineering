package model

type Chapter struct {
	ID      int64  `json:"id"`
	BookID  int64  `json:"bookid"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	Price   int64  `json:"price"`
	Content string `json:"content"`
}
