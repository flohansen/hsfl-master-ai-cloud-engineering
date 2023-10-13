package model

type Book struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
}
